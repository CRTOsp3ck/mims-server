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

// Database settings - development
const (
	host     = "localhost"
	port     = 5432 // Default port
	user     = "sp3ck"
	password = "88888888"
	dbname   = "mims_server_development"
)

// Database settings - production
// const (
// 	host     = "127.0.0.1"
// 	port     = 5432 // Default port
// 	user     = "root"
// 	password = ""
// 	dbname   = "root"
// )

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
	app.Get("/ag", func(c *fiber.Ctx) error {
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
	app.Get("/ag/:user", func(c *fiber.Ctx) error {
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
	app.Get("/op", func(c *fiber.Ctx) error {
		rows, err := db.Query("SELECT id, start_time, end_time, location, agent_id, total_sales_qty, total_cost, total_sales_amount, net_profit, balance_id, inventory_id, created_at, updated_at FROM operation ORDER BY id")
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		defer rows.Close()
		result := Operations{}

		for rows.Next() {
			res := Operation{}
			err := rows.Scan(&res.ID, &res.StartTime, &res.EndTime, &res.Location, &res.AgentID, &res.TotalSalesQty, &res.TotalCost, &res.TotalSalesAmount, &res.NetProfit, &res.BalanceID, &res.InventoryID, &res.CreatedAt, &res.UpdatedAt)
			if err != nil {
				return err
			}
			result.Operations = append(result.Operations, res)
		}

		return c.JSON(result)
	})
	// Start operation
	app.Post("/op/start/:location-:agent_user/bal/:start_bal_cash-:start_bal_qr/inv/:start_item_bal", func(c *fiber.Ctx) error {
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
		sbItem := c.Params("start_bal")
		res, err = db.Query("INSERT INTO inventory (start_bal, end_bal, created_at, updated_at)VALUES ($1, $2, $3, $4)", sbItem, "-1", time.Now(), time.Now())
		_ = res
		if err != nil {
			return err
		}
		inv := new(Inventory)
		// Re-querying because the scan from insert has no value?
		resReQuery = db.QueryRow("SELECT id, start_bal, end_bal, created_at, updated_at FROM inventory ORDER BY ID DESC LIMIT 1")
		resReQuery.Scan(&inv.ID, &inv.StartBal, &inv.EndBal, &inv.CreatedAt, &inv.UpdatedAt)
		log.Println("Inventory id", inv.ID)
		paramCache.InventoryID = inv.ID

		// Insert all cached data to db
		res, err = db.Query("INSERT INTO operation (start_time, end_time, location, agent_id, total_sales_qty, total_cost, total_sales_amount, net_profit, balance_id, inventory_id, created_at, updated_at)VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
			paramCache.StartTime, paramCache.EndTime, paramCache.Location, paramCache.AgentID, -1, -1.00, -1.00, -1.00, paramCache.BalanceID, paramCache.InventoryID, time.Now(), time.Now())
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
	app.Put("/op/end/:operation_id-:total_cost", func(c *fiber.Ctx) error {
		paramCache := new(Operation)
		paramCache.EndTime = time.Now()

		// Find and calculate all sales from this operation (using operation_id)
		opid, err := strconv.Atoi(c.Params("operation_id"))
		if err != nil {
			return err
		}
		rows, err := db.Query("SELECT id, amount, quantity, payment_type, operation_id, item_id, created_at, updated_at FROM sale WHERE operation_id=$1", opid)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		defer rows.Close()
		sales := Sales{}
		for rows.Next() {
			sale := Sale{}
			if err := rows.Scan(&sale.ID, &sale.Amount, &sale.Qty, &sale.PaymentType, &sale.OperationID, &sale.ItemID, &sale.CreatedAt, &sale.UpdatedAt); err != nil {
				return err // Exit if we get an error
			}

			// Append Sale to Sales
			sales.Sales = append(sales.Sales, sale)
		}

		totalQty := 0.00
		totalSales := 0.00
		for _, sa := range sales.Sales {
			totalQty += float64(sa.Qty)
			totalSales += float64(sa.Amount)
		}
		paramCache.TotalSalesQty = int(totalQty)

		// Enter total cost during operation end
		totalCost, _ := strconv.ParseFloat(c.Params("total_cost"), 32)
		paramCache.TotalCost = float32(totalCost)

		// Calculate total sales amount (sale qty*price sold)
		paramCache.TotalSalesAmount = float32(totalSales)

		// Calculate net profit (total sales qty * rm8)
		paramCache.NetProfit = float32(totalSales - totalCost)
		paramCache.UpdatedAt = time.Now()

		// Update operation into database
		res, err := db.Query("UPDATE operation SET end_time=$1,total_sales_qty=$2,total_cost=$3,total_sales_amount=$4,net_profit=$5, updated_at=$6 WHERE id=$7", paramCache.EndTime, paramCache.TotalSalesQty, paramCache.TotalCost, paramCache.TotalSalesAmount, paramCache.NetProfit, paramCache.UpdatedAt, opid)
		_ = res
		if err != nil {
			return err
		}

		op := Operation{}
		// Re-querying because the scan from insert has no value?
		resReQuery := db.QueryRow("SELECT id, start_time, end_time, location, agent_id, total_sales_qty, total_cost, total_sales_amount, net_profit, balance_id, inventory_id, created_at, updated_at FROM operation WHERE id=$1", opid)
		resReQuery.Scan(&op.ID, &op.StartTime, &op.EndTime, &op.Location, &op.AgentID, &op.TotalSalesQty, &op.TotalCost, &op.TotalSalesAmount, &op.NetProfit, &op.BalanceID, &op.InventoryID, &op.CreatedAt, &op.UpdatedAt)

		// Return operation in JSON format
		return c.JSON(op)

		// TODO: Need to update the balances and inventories also!! -> work on logging in
	})

	// >> Sales
	// Get all sales
	app.Get("/sa", func(c *fiber.Ctx) error {
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
	app.Get("/sa/find/:syear-:smonth-:sday-:eyear-:emonth-:eday", func(c *fiber.Ctx) error {
		//argh hungry. eat first la. go to vending machine first. kk.
		//get some gummy bearsx2, melon milk, yougrt drink, soy drink
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
	// Get list of sale (by operation_id)
	app.Get("/sa/find/:operation_id", func(c *fiber.Ctx) error {
		// i really should eat, man. im losing weight. you want to get gastric again now, do you?
		opid := c.Params("operation_id")

		rows, err := db.Query("SELECT id, amount, quantity, payment_type, operation_id, item_id, created_at, updated_at FROM sale WHERE operation_id=$1", opid)
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
	app.Post("/sa/new/:amount-:qty-:payment_type-:operation_id-:item_id-:group_sale_id?", func(c *fiber.Ctx) error {
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

		// if gsid is not present, the default value 0 means its a parent
		if c.Params("group_sale_id") != "" {
			gsid, err := strconv.Atoi(c.Params("group_sale_id"))
			if err != nil {
				return err
			}
			paramCache.GroupSaleID = gsid
		}

		// Insert sale into database
		res, err := db.Query("INSERT INTO sale (amount, quantity, payment_type, operation_id, item_id, group_sale_id, created_at, updated_at)VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
			paramCache.Amount, paramCache.Qty, paramCache.PaymentType, paramCache.OperationID, paramCache.ItemID, paramCache.GroupSaleID, time.Now(), time.Now())
		_ = res
		if err != nil {
			return err
		}

		sale := new(Sale)
		// Re-querying because the scan from insert has no value?
		resReQuery := db.QueryRow("SELECT id, amount, quantity, payment_type, operation_id, item_id, group_sale_id, created_at, updated_at FROM sale ORDER BY ID DESC LIMIT 1")
		resReQuery.Scan(&sale.ID, &sale.Amount, &sale.Qty, &sale.PaymentType, &sale.OperationID, &sale.ItemID, &sale.GroupSaleID, &sale.CreatedAt, &sale.UpdatedAt)

		// Return Employee in JSON format
		return c.JSON(sale)
	})
	// Update Sale

	// Delete Sale (admin only / only the most recent one)

	// >> Inventory
	// Get all inventory
	app.Get("/inv", func(c *fiber.Ctx) error {
		rows, err := db.Query("SELECT id, start_bal, end_bal, created_at, updated_at FROM inventory order by id")
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		defer rows.Close()
		result := Inventories{}

		for rows.Next() {
			inv := Inventory{}
			if err := rows.Scan(&inv.ID, &inv.StartBal, &inv.EndBal, &inv.CreatedAt, &inv.UpdatedAt); err != nil {
				return err // Exit if we get an error
			}

			// Append Sale to Sales
			result.Inventories = append(result.Inventories, inv)
		}
		// Return Sales in JSON format
		return c.JSON(result)
	})
	// New inventory
	app.Post("/inv/new/:start_bal", func(c *fiber.Ctx) error {
		// Insert a new inventory record into database
		sbItem := c.Params("start_bal")
		res, err := db.Query("INSERT INTO inventory (start_bal, end_bal, created_at, updated_at)VALUES ($1, $2, $3, $4)", sbItem, "-1", time.Now(), time.Now())
		_ = res
		if err != nil {
			return err
		}
		inv := new(Inventory)
		// Re-querying because the scan from insert has no value?
		resReQuery := db.QueryRow("SELECT id, start_bal, end_bal, created_at, updated_at FROM inventory ORDER BY ID DESC LIMIT 1")
		resReQuery.Scan(&inv.ID, &inv.StartBal, &inv.EndBal, &inv.CreatedAt, &inv.UpdatedAt)

		// Return Inventory in JSON format
		return c.JSON(inv)
	})
	// Update inventory
	app.Put("/inv/up/:id-:end_bal", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return err
		}
		eb := c.Params("end_bal")

		// Update inventory into database
		res, err := db.Query("UPDATE inventory SET end_bal=$1,updated_at=$2 WHERE id=$3", eb, time.Now(), id)
		_ = res
		if err != nil {
			return err
		}

		inv := Inventory{}
		// Re-querying because the scan from insert has no value?
		resReQuery := db.QueryRow("SELECT id, start_bal, end_bal, created_at, updated_at FROM inventory ORDER BY ID DESC LIMIT 1")
		resReQuery.Scan(&inv.ID, &inv.StartBal, &inv.EndBal, &inv.CreatedAt, &inv.UpdatedAt)

		// Return inv in JSON format
		return c.JSON(inv)
	})
	// Delete inventory (admin only)

	// >> Item
	// Get all items
	app.Get("/it", func(c *fiber.Ctx) error {
		rows, err := db.Query("SELECT id, name, des, price, cost_price, min_combo_qty, min_combo_price, created_at, updated_at FROM item order by id")
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		defer rows.Close()
		result := Items{}

		for rows.Next() {
			it := Item{}
			if err := rows.Scan(&it.ID, &it.Name, &it.Des, &it.Price, &it.Cost, &it.MinComboQty, &it.MinComboPrice, &it.CreatedAt, &it.UpdatedAt); err != nil {
				return err // Exit if we get an error
			}

			// Append Sale to Sales
			result.Items = append(result.Items, it)
		}
		// Return Sales in JSON format
		return c.JSON(result)
	})

	log.Fatal(app.Listen(":3001"))
}

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

type Inventories struct {
	Inventories []Inventory `json:"inventories"`
}

type Inventory struct {
	ID        int       `json:"id"`
	StartBal  string    `json:"start_bal"`
	EndBal    string    `json:"end_bal"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Sale struct {
	ID          int       `json:"id"`
	Amount      float32   `json:"amount"`
	Qty         float32   `json:"quantity"` //**(SHOULD CHANGE TO INT) this is float and not int bcos in case we plan to sell by weight, then it wouldnt make sense to use int
	PaymentType int       `json:"payment_type"`
	OperationID int       `json:"operation_id"`
	ItemID      int       `json:"item_id"`
	GroupSaleID int       `json:"group_sale_id"`
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

type Operations struct {
	Operations []Operation `json:"operations"`
}
type Item struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Des           string    `json:"des"`
	Price         float32   `json:"price"`
	Cost          float32   `json:"cost_price"`
	MinComboQty   int       `json:"min_combo_qty"`
	MinComboPrice float32   `json:"min_combo_price"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Items struct {
	Items []Item `json:"items"`
}
