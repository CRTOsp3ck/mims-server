// 🚀 Fiber is an Express inspired web framework written in Go with 💖
// 📌 API Documentation: https://docs.gofiber.io
// 📝 Github Repository: https://github.com/gofiber/fiber
package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

// Database instance
var db *sql.DB

// Database settings
const (
	host     = "localhost"
	port     = 5432 // Default port
	user     = "sp3ck"
	password = "88888888"
	dbname   = "mims_server_development"
)

// Data models
type Agent struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	IsOwner   bool      `json:"is_owner"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Agents struct {
	Agents []Agent `json:"agents"`
}

type Balance struct {
	ID        int       `json:"id"`
	BalCash   string    `json:"bal_cash"`
	BalQr     string    `json:"bal_qr"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Inventory struct {
	ID           int       `json:"id"`
	StartItemBal string    `json:"start_item_bal"`
	EndItemBal   string    `json:"end_item_bal"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Sale struct {
	ID          int       `json:"id"`
	Amount      float32   `json:"amount"`
	Qty         float32   `json:"quantity"` //this is float and not int bcos in case we plan to sell by weight, then it wouldnt make sense to use int
	PaymentType int       `json:"payment_type"`
	OperationID int       `json:"operation_id"`
	ItemID      int       `json:"item_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Sales struct {
	Sales []Sale `json:"sales"`
}

type Operation struct {
	ID               int       `json:"id"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	Location         string    `json:"location"`
	AgentID          int       `json:"agent_id"`
	TotalSalesQty    int       `json:"total_sales_qty"`
	TotalCost        float32   `json:"total_cost"`
	TotalSalesAmount float32   `json:"total_sales_amount"`
	NetProfit        float32   `json:"net_profit"`
	BalanceID        int       `json:"balance_id"`
	InventoryID      int       `json:"inventory_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type Item struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Des       string    `json:"des"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Connect function
func Connect() error {
	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname))
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}
	return nil
}

func main() {
	// Connect with database
	if err := Connect(); err != nil {
		log.Fatal(err)
	}

	// Create a Fiber app
	app := fiber.New()

	// >> Agent

	// Get list of all agents
	app.Get("/agents", func(c *fiber.Ctx) error {
		rows, err := db.Query("SELECT id, username, password, name, email, phone, is_owner, created_at, updated_at FROM agent order by id")
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		defer rows.Close()
		result := Agents{}

		for rows.Next() {
			res := Agent{}
			err := rows.Scan(&res.ID, &res.Username, &res.Password, &res.Name, &res.Email, &res.Phone, &res.IsOwner, &res.CreatedAt, &res.UpdatedAt)
			if err != nil {
				return err
			}
			result.Agents = append(result.Agents, res)
		}

		return c.JSON(result)
	})

	// Find agent by username
	app.Get("/agents/:user", func(c *fiber.Ctx) error {
		user := c.Params("user")
		row := db.QueryRow("SELECT id, username, password, name, email, phone, is_owner, created_at, updated_at FROM agent WHERE username = $1", user)
		res := Agent{}
		err := row.Scan(&res.ID, &res.Username, &res.Password, &res.Name, &res.Email, &res.Phone, &res.IsOwner, &res.CreatedAt, &res.UpdatedAt)
		if err != nil {
			return err
		}
		return c.JSON(res)
	})

	// >> Operation

	// Get all operations

	// Start operation
	app.Post("/operations/start/:location-:agent_user/bal/:start_bal_cash-:start_bal_qr/inv/:start_item_bal",
		func(c *fiber.Ctx) error {
			paramCache := new(Operation)
			paramCache.StartTime = time.Now()
			paramCache.EndTime = time.Time{}
			paramCache.Location = c.Params("location")

			// Find ID of the agent's username
			user := c.Params("agent_user")
			row := db.QueryRow("SELECT id, username, password, name, email, phone, is_owner, created_at, updated_at FROM agent WHERE username = $1", user)
			agent := Agent{}
			err := row.Scan(&agent.ID, &agent.Username, &agent.Password, &agent.Name, &agent.Email, &agent.Phone, &agent.IsOwner, &agent.CreatedAt, &agent.UpdatedAt)
			if err != nil {
				return err
			}
			log.Println("Agent id", agent.ID)
			paramCache.AgentID = agent.ID

			// Insert a new balance record into database
			sbCash := c.Params("start_bal_cash")
			sbQr := c.Params("start_bal_qr")
			balCashStr := "sb=" + sbCash + "&eb=-1" //-1 means operation is in progress
			balQrStr := "sb=" + sbQr + "&eb=-1"
			res, err := db.Query("INSERT INTO balance (bal_cash, bal_qr, created_at, updated_at)VALUES ($1, $2, $3, $4)", balCashStr, balQrStr, time.Now(), time.Now())
			_ = res
			if err != nil {
				return err
			}
			bal := new(Balance)
			// Re-querying because the scan from insert has no value?
			resReQuery := db.QueryRow("SELECT id, bal_cash, bal_qr, created_at, updated_at FROM balance ORDER BY ID DESC LIMIT 1")
			resReQuery.Scan(&bal.ID, &bal.BalCash, &bal.BalQr, &bal.CreatedAt, &bal.UpdatedAt)
			log.Println("Balance id", bal.ID)
			paramCache.BalanceID = bal.ID

			// Insert a new inventory record into database
			sbItem := c.Params("start_item_bal")
			res, err = db.Query("INSERT INTO inventory (start_item_bal, end_item_bal, created_at, updated_at)VALUES ($1, $2, $3, $4)", sbItem, "", time.Now(), time.Now())
			_ = res
			if err != nil {
				return err
			}
			inv := new(Inventory)
			// Re-querying because the scan from insert has no value?
			resReQuery = db.QueryRow("SELECT id, start_item_bal, end_item_bal, created_at, updated_at FROM inventory ORDER BY ID DESC LIMIT 1")
			resReQuery.Scan(&inv.ID, &inv.StartItemBal, &inv.EndItemBal, &inv.CreatedAt, &inv.UpdatedAt)
			log.Println("Inventory id", inv.ID)
			paramCache.InventoryID = inv.ID

			// Insert all cached data to db
			res, err = db.Query("INSERT INTO operation (start_time, end_time, location, agent_id, total_sales_qty, total_cost, total_sales_amount, net_profit, balance_id, inventory_id, created_at, updated_at)VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
				paramCache.StartTime, paramCache.EndTime, paramCache.Location, paramCache.AgentID, 0, 0.00, 0.00, 0.00, paramCache.BalanceID, paramCache.InventoryID, time.Now(), time.Now())
			_ = res
			if err != nil {
				return err
			}
			op := new(Operation)
			// Re-querying because the scan from insert has no value?
			resReQuery = db.QueryRow("SELECT id, start_time, end_time, location, agent_id, total_sales_qty, total_cost, total_sales_amount, net_profit, balance_id, inventory_id, created_at, updated_at FROM operation ORDER BY ID DESC LIMIT 1")
			resReQuery.Scan(&op.ID, &op.StartTime, &op.EndTime, &op.Location, &op.AgentID, &op.TotalSalesQty, &op.TotalCost, &op.TotalSalesAmount, &op.NetProfit, &op.BalanceID, &op.InventoryID, &op.CreatedAt, &op.UpdatedAt)

			// Return final operation in JSON format
			return c.JSON(op)
		})

	// End operation
	app.Put("/operations/end/:id",
		func(c *fiber.Ctx) error {
			id, err := strconv.Atoi(c.Params("id"))
			if err != nil {
				return err
			}

			paramCache := new(Operation)
			paramCache.EndTime = time.Now()
			// Find and calculate all sales from this operation (using oepration_id)
			paramCache.TotalSalesQty = 0
			// Enter total cost during operation end
			paramCache.TotalCost = 0.00
			// Calculate total sales amount (sale qty*price sold)
			paramCache.TotalSalesAmount = 0.00
			// Calculate net profit (total sales qty * rm8)
			paramCache.NetProfit = 0.00
			paramCache.UpdatedAt = time.Now()

			// Update operation into database
			res, err := db.Query("UPDATE operation SET end_time=$1,total_sales_qty=$2,total_cost=$3,total_sales_amount=$4,net_profit=$5, updated_at=$6 WHERE id=$7", paramCache.EndTime, paramCache.TotalSalesQty, paramCache.TotalCost, paramCache.TotalSalesAmount, paramCache.NetProfit, paramCache.UpdatedAt, id)
			_ = res
			if err != nil {
				return err
			}

			op := Operation{}
			// Re-querying because the scan from insert has no value?
			resReQuery := db.QueryRow("SELECT id, start_time, end_time, location, agent_id, total_sales_qty, total_cost, total_sales_amount, net_profit, balance_id, inventory_id, created_at, updated_at FROM operation ORDER BY ID DESC LIMIT 1")
			resReQuery.Scan(&op.ID, &op.StartTime, &op.EndTime, &op.Location, &op.AgentID, &op.TotalSalesQty, &op.TotalCost, &op.TotalSalesAmount, &op.NetProfit, &op.BalanceID, &op.InventoryID, &op.CreatedAt, &op.UpdatedAt)

			// Return operation in JSON format
			return c.Status(201).JSON(op)
		})

	// >> Sales

	// Get all sales
	app.Get("/sales/find/", func(c *fiber.Ctx) error {
		rows, err := db.Query("SELECT id, amount, quantity, payment_type, operation_id, item_id, created_at, updated_at FROM sale order by id")
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		defer rows.Close()
		result := Sales{}

		for rows.Next() {
			sale := Sale{}
			if err := rows.Scan(&sale.ID, &sale.Amount, &sale.Qty, &sale.PaymentType, &sale.OperationID, &sale.ItemID, &sale.CreatedAt, &sale.UpdatedAt); err != nil {
				return err // Exit if we get an error
			}

			// Append Sale to Sales
			result.Sales = append(result.Sales, sale)
		}
		// Return Sales in JSON format
		return c.JSON(result)
	})

	// Get list of sale (by date)
	app.Get("/sales/find/:syear-:smonth-:sday-:eyear-:emonth-:eday", func(c *fiber.Ctx) error {

		s_year, _ := strconv.Atoi(c.Params("syear"))
		s_month, _ := strconv.Atoi(c.Params("smonth"))
		s_day, _ := strconv.Atoi(c.Params("sday"))

		e_year, _ := strconv.Atoi(c.Params("eyear"))
		e_month, _ := strconv.Atoi(c.Params("emonth"))
		e_day, _ := strconv.Atoi(c.Params("eday"))

		sd := time.Date(s_year, time.Month(s_month), s_day, 0, 0, 0, 0, time.Local)
		ed := time.Date(e_year, time.Month(e_month), e_day, 0, 0, 0, 0, time.Local)

		rows, err := db.Query("SELECT id, amount, quantity, payment_type, operation_id, item_id, created_at, updated_at FROM sale WHERE date_trunc('day', created_at) BETWEEN $1 and $2", sd, ed)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		defer rows.Close()
		result := Sales{}

		for rows.Next() {
			sale := Sale{}
			if err := rows.Scan(&sale.ID, &sale.Amount, &sale.Qty, &sale.PaymentType, &sale.OperationID, &sale.ItemID, &sale.CreatedAt, &sale.UpdatedAt); err != nil {
				return err // Exit if we get an error
			}

			// Append Sale to Sales
			result.Sales = append(result.Sales, sale)
		}
		// Return Sales in JSON format
		return c.JSON(result)
	})

	// New Sale
	app.Post("/sales/new/:amount-:qty-:payment_type-:operation_id-:item_id", func(c *fiber.Ctx) error {
		paramCache := new(Sale)

		amt, err := strconv.Atoi(c.Params("amount"))
		if err != nil {
			return err
		}
		paramCache.Amount = float32(amt)

		qty, err := strconv.Atoi(c.Params("qty"))
		if err != nil {
			return err
		}
		paramCache.Qty = float32(qty)

		pt, err := strconv.Atoi(c.Params("payment_type"))
		if err != nil {
			return err
		}
		paramCache.PaymentType = pt

		oid, err := strconv.Atoi(c.Params("operation_id"))
		if err != nil {
			return err
		}
		paramCache.OperationID = oid

		iid, err := strconv.Atoi(c.Params("item_id"))
		if err != nil {
			return err
		}
		paramCache.ItemID = iid

		// Insert sale into database
		res, err := db.Query("INSERT INTO sale (amount, quantity, payment_type, operation_id, item_id, created_at, updated_at)VALUES ($1, $2, $3, $4, $5, $6, $7)",
			paramCache.Amount, paramCache.Qty, paramCache.PaymentType, paramCache.OperationID, paramCache.ItemID, time.Now(), time.Now())
		_ = res
		if err != nil {
			return err
		}

		sale := new(Sale)
		// Re-querying because the scan from insert has no value?
		resReQuery := db.QueryRow("SELECT id, amount, quantity, payment_type, operation_id, item_id, created_at, updated_at FROM sale ORDER BY ID DESC LIMIT 1")
		resReQuery.Scan(&sale.ID, &sale.Amount, &sale.Qty, &sale.PaymentType, &sale.OperationID, &sale.ItemID, &sale.CreatedAt, &sale.UpdatedAt)

		// Print result
		log.Println(sale)

		// Return Employee in JSON format
		return c.JSON(sale)
	})

	// Update Sale

	// Delete Sale (admin only / only the most recent one)

	// >> Inventory

	// Add inventory

	// Update inventory

	// Delete inventory (admin only)

	log.Fatal(app.Listen(":3000"))
}
