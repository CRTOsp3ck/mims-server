// 🚀 Fiber is an Express inspired web framework written in Go with 💖
// 📌 API Documentation: https://docs.gofiber.io
// 📝 Github Repository: https://github.com/gofiber/fiber
package main

import (
	"database/sql"
	"fmt"
	"log"
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
	Qty         float32   `json:"quantity"`
	PaymentType int       `json:"payment_type"`
	OperationID int       `json:"operation_id"`
	ItemID      int       `json:"item_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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

	// Start operation

	// End operation

	// >> Sales

	// New Sale

	// Update Sale

	// Delete Sale (admin only / only the most recent one)

	// >> Inventory

	// Add inventory

	// Update inventory

	// Delete inventory (admin only)

	log.Fatal(app.Listen(":3000"))
}
