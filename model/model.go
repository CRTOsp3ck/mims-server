package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"password"`
	Names    string `json:"names"`
}

type Product struct {
	gorm.Model
	Title       string `gorm:"not null" json:"title"`
	Description string `gorm:"not null" json:"description"`
	Amount      int    `gorm:"not null" json:"amount"`
}

type Agent struct {
	gorm.Model
	Username string `gorm:"not null" json:"username"`
	Password string `gorm:"not null" json:"password"`
	Name     string `gorm:"not null" json:"name"`
	Email    string `gorm:"not null" json:"email"`
	Phone    string `gorm:"not null" json:"phone"`
	IsOwner  bool   `gorm:"not null" json:"is_owner"`
}

type Balance struct {
	gorm.Model
	BalCash string `gorm:"not null" json:"bal_cash"`
	BalQr   string `gorm:"not null" json:"bal_qr"`
}

type Inventory struct {
	gorm.Model
	StartBal string `gorm:"not null" json:"start_bal"`
	EndBal   string `gorm:"not null" json:"end_bal"`
}

type Sale struct {
	gorm.Model
	Amount      float32 `gorm:"not null" json:"amount"`
	Qty         int     `gorm:"not null" json:"qty"`
	PaymentType int     `gorm:"not null" json:"payment_type"`
	OperationID int     `gorm:"not null" json:"operation_id"`
	ItemID      int     `gorm:"not null" json:"item_id"`
	GroupSaleID int     `gorm:"not null" json:"group_sale_id"`
}

type Operation struct {
	gorm.Model
	StartTime        time.Time `gorm:"not null" json:"start_time"`
	EndTime          time.Time `gorm:"not null" json:"end_time"`
	Location         string    `gorm:"not null" json:"location"`
	AgentID          int       `gorm:"not null" json:"agent_id"`
	TotalSalesQty    int       `gorm:"not null" json:"total_sales_qty"`
	TotalCost        float32   `gorm:"not null" json:"total_cost"`
	TotalSalesAmount float32   `gorm:"not null" json:"total_sales_amount"`
	NetProfit        float32   `gorm:"not null" json:"net_profit"`
	BalanceID        int       `gorm:"not null" json:"balance_id"`
	InventoryID      int       `gorm:"not null" json:"inventory_id"`
}

type Item struct {
	gorm.Model
	Name          string  `gorm:"not null" json:"name"`
	Des           string  `gorm:"not null" json:"des"`
	Price         float32 `gorm:"not null" json:"price"`
	Cost          float32 `gorm:"not null" json:"cost_price"`
	MinComboQty   int     `gorm:"not null" json:"min_combo_qty"`
	MinComboPrice float32 `gorm:"not null" json:"min_combo_price"`
}
