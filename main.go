package main

import (
	"api-fiber-gorm/database"
	"api-fiber-gorm/model"
	"api-fiber-gorm/router"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())

	database.ConnectDB()

	router.SetupRoutes(app)

	// >> Agent
	// Get list of all agents
	app.Get("/ag", func(c *fiber.Ctx) error {
		db := database.DB
		var ag []model.Agent
		db.Find(&ag)
		return c.JSON(ag)
	})
	// Find agent by username
	app.Get("/ag/:user", func(c *fiber.Ctx) error {
		db := database.DB
		var ag model.Agent
		db.First(&ag, "username = ?", c.Params("user"))
		return c.JSON(ag)
	})

	// >> Operation
	// Get all operations
	app.Get("/op", func(c *fiber.Ctx) error {
		db := database.DB
		var op []model.Operation
		db.Find(&op)
		return c.JSON(op)
	})
	// Start new operation
	app.Post("/op/start/:location-:agent_user/bal/:start_bal_cash-:start_bal_qr/inv/:start_item_bal", func(c *fiber.Ctx) error {
		db := database.DB

		op := new(model.Operation)
		op.StartTime = time.Now()
		op.EndTime = time.Time{}
		op.Location = c.Params("location")

		// Find ID of the agent's username
		var ag model.Agent
		db.First(&ag, "username = ?", c.Params("agent_user"))
		op.AgentID = int(ag.ID)

		// Insert a new balance record into database
		sbCash := c.Params("start_bal_cash")
		sbQr := c.Params("start_bal_qr")

		var bal model.Balance
		bal.BalCash = "sb=" + sbCash + "&eb=-1" //-1 means operation is in progress
		bal.BalQr = "sb=" + sbQr + "&eb=-1"
		db.Create(&bal)
		op.BalanceID = int(bal.ID)

		// Insert a new inventory record into database
		var inv model.Inventory
		inv.StartBal = c.Params("start_item_bal")
		inv.EndBal = "-1"
		db.Create(&inv)
		op.InventoryID = int(inv.ID)

		// Insert all cached data to db
		db.Create(&op)

		// Return final operation in JSON format
		return c.JSON(op)
	})
	// End operation
	app.Put("/op/end/:op_id-:total_cost", func(c *fiber.Ctx) error {
		db := database.DB
		opid := c.Params("op_id")

		// find row with this operation_id in db
		var op model.Operation
		db.Where("id = ?", opid).Find(&op)
		op.EndTime = time.Now()

		// Find and calculate all sales from this operation (using operation_id)
		var sales []model.Sale
		db.Where("operation_id = ?", opid).Find(&sales)

		totalQty := 0.00
		totalSales := 0.00
		for _, sa := range sales {
			totalQty += float64(sa.Qty)
			totalSales += float64(sa.Amount)
		}
		op.TotalSalesQty = int(totalQty)

		// Enter total cost during operation end
		totalCost, _ := strconv.ParseFloat(c.Params("total_cost"), 32)
		op.TotalCost = float32(totalCost)

		// Calculate total sales amount (sale qty*price sold)
		op.TotalSalesAmount = float32(totalSales)

		// Calculate net profit (total sales qty * rm8)
		op.NetProfit = float32(totalSales - totalCost)

		// Update operation into database
		db.Save(&op)

		// Return operation in JSON format
		return c.JSON(op)

		// TODO: Need to update the balances and inventories also!! -> work on logging in
	})

	// >> Sales
	// Get all sales
	app.Get("/sa", func(c *fiber.Ctx) error {
		db := database.DB
		var sales []model.Sale
		db.Find(&sales)
		// Return Sales in JSON format
		return c.JSON(sales)
	})
	// Get list of sale (by date)
	app.Get("/sa/find/:syear-:smonth-:sday-:eyear-:emonth-:eday", func(c *fiber.Ctx) error {
		s_year, _ := strconv.Atoi(c.Params("syear"))
		s_month, _ := strconv.Atoi(c.Params("smonth"))
		s_day, _ := strconv.Atoi(c.Params("sday"))

		e_year, _ := strconv.Atoi(c.Params("eyear"))
		e_month, _ := strconv.Atoi(c.Params("emonth"))
		e_day, _ := strconv.Atoi(c.Params("eday"))

		sd := time.Date(s_year, time.Month(s_month), s_day, 0, 0, 0, 0, time.Local) //should add the time also
		ed := time.Date(e_year, time.Month(e_month), e_day, 0, 0, 0, 0, time.Local)

		db := database.DB
		var sales []model.Sale
		db.Where("created_at BETWEEN ? AND ?", sd, ed).Find(&sales)

		// Return Sales in JSON format
		return c.JSON(sales)
	})
	// Get list of sale (by operation_id)
	app.Get("/sa/find/:op_id", func(c *fiber.Ctx) error {
		opid, _ := strconv.Atoi(c.Params("op_id"))

		db := database.DB
		var sales []model.Sale
		db.Where("operation_id <> ?", opid-1).Find(&sales)

		// Return Sales in JSON format
		return c.JSON(sales)

	})
	// New Sale
	app.Post("/sa/new/:amount-:qty-:payment_type-:op_id-:item_id-:group_sale_id", func(c *fiber.Ctx) error {
		db := database.DB
		var sale model.Sale
		amt, err := strconv.Atoi(c.Params("amount"))
		if err != nil {
			return err
		}
		sale.Amount = float32(amt)

		qty, err := strconv.Atoi(c.Params("qty"))
		if err != nil {
			return err
		}
		sale.Qty = qty

		pt, err := strconv.Atoi(c.Params("payment_type"))
		if err != nil {
			return err
		}
		sale.PaymentType = pt

		oid, err := strconv.Atoi(c.Params("op_id"))
		if err != nil {
			return err
		}
		sale.OperationID = oid

		iid, err := strconv.Atoi(c.Params("item_id"))
		if err != nil {
			return err
		}
		sale.ItemID = iid

		// if gsid is not present, the default value 0 means its a parent
		if c.Params("group_sale_id") != "" {
			gsid, err := strconv.Atoi(c.Params("group_sale_id"))
			if err != nil {
				return err
			}
			sale.GroupSaleID = gsid
		}

		db.Create(&sale)

		// Return sale in JSON format
		return c.JSON(sale)
	})

	// Delete Sale (admin only / only the most recent one)
	app.Post("/sa/del/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		db := database.DB

		var sale model.Sale
		db.First(&sale, id)
		db.Delete(&sale)
		return c.JSON(sale)
	})

	// // >> Inventory
	// // Get all inventory
	// app.Get("/inv", func(c *fiber.Ctx) error {
	// 	rows, err := db.Query("SELECT id, start_bal, end_bal, created_at, updated_at FROM inventory order by id")
	// 	if err != nil {
	// 		return c.Status(500).SendString(err.Error())
	// 	}
	// 	defer rows.Close()
	// 	result := Inventories{}

	// 	for rows.Next() {
	// 		inv := Inventory{}
	// 		if err := rows.Scan(&inv.ID, &inv.StartBal, &inv.EndBal, &inv.CreatedAt, &inv.UpdatedAt); err != nil {
	// 			return err // Exit if we get an error
	// 		}

	// 		// Append Sale to Sales
	// 		result.Inventories = append(result.Inventories, inv)
	// 	}
	// 	// Return Sales in JSON format
	// 	return c.JSON(result)
	// })
	// // New inventory
	// app.Post("/inv/new/:start_bal", func(c *fiber.Ctx) error {
	// 	// Insert a new inventory record into database
	// 	sbItem := c.Params("start_bal")
	// 	res, err := db.Query("INSERT INTO inventory (start_bal, end_bal, created_at, updated_at)VALUES ($1, $2, $3, $4)", sbItem, "-1", time.Now(), time.Now())
	// 	_ = res
	// 	if err != nil {
	// 		return err
	// 	}
	// 	inv := new(Inventory)
	// 	// Re-querying because the scan from insert has no value?
	// 	resReQuery := db.QueryRow("SELECT id, start_bal, end_bal, created_at, updated_at FROM inventory ORDER BY ID DESC LIMIT 1")
	// 	resReQuery.Scan(&inv.ID, &inv.StartBal, &inv.EndBal, &inv.CreatedAt, &inv.UpdatedAt)

	// 	// Return Inventory in JSON format
	// 	return c.JSON(inv)
	// })
	// // Update inventory
	// app.Put("/inv/up/:id-:end_bal", func(c *fiber.Ctx) error {
	// 	id, err := strconv.Atoi(c.Params("id"))
	// 	if err != nil {
	// 		return err
	// 	}
	// 	eb := c.Params("end_bal")

	// 	// Update inventory into database
	// 	res, err := db.Query("UPDATE inventory SET end_bal=$1,updated_at=$2 WHERE id=$3", eb, time.Now(), id)
	// 	_ = res
	// 	if err != nil {
	// 		return err
	// 	}

	// 	inv := Inventory{}
	// 	// Re-querying because the scan from insert has no value?
	// 	resReQuery := db.QueryRow("SELECT id, start_bal, end_bal, created_at, updated_at FROM inventory ORDER BY ID DESC LIMIT 1")
	// 	resReQuery.Scan(&inv.ID, &inv.StartBal, &inv.EndBal, &inv.CreatedAt, &inv.UpdatedAt)

	// 	// Return inv in JSON format
	// 	return c.JSON(inv)
	// })
	// // Delete inventory (admin only)

	// // >> Item
	// // Get all items
	// app.Get("/it", func(c *fiber.Ctx) error {
	// 	rows, err := db.Query("SELECT id, name, des, price, cost_price, min_combo_qty, min_combo_price, created_at, updated_at FROM item order by id")
	// 	if err != nil {
	// 		return c.Status(500).SendString(err.Error())
	// 	}
	// 	defer rows.Close()
	// 	result := Items{}

	// 	for rows.Next() {
	// 		it := Item{}
	// 		if err := rows.Scan(&it.ID, &it.Name, &it.Des, &it.Price, &it.Cost, &it.MinComboQty, &it.MinComboPrice, &it.CreatedAt, &it.UpdatedAt); err != nil {
	// 			return err // Exit if we get an error
	// 		}

	// 		// Append Sale to Sales
	// 		result.Items = append(result.Items, it)
	// 	}
	// 	// Return Sales in JSON format
	// 	return c.JSON(result)
	// })

	log.Fatal(app.Listen(":3001"))
}
