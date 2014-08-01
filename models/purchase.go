package models

import (
	// External Packages
	. "github.com/pfacheris/kickback/db"
	"time"
)

type Purchase struct {
	Id                 int64     `json:"id"`
	PurchasePrice      float32   `json:"purchase_price" binding:"required" sql:"type:decimal(11,2);not null"`
	KickbackAmount     float32   `json:"kickback_amount" binding:"required" sql:"type:decimal(11,2);not null"`
	PurchaseAt         time.Time `json:"purchase_at"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	DeletedAt          time.Time `json:"-"`
	UserId             int64     `json:"-"`
	ProductId          int64     `json:"-"`
	SellerName         string    `json:"seller_name"`
	CurrentSellerPrice float32   `json:"current_price" binding:"required" sql:"type:decimal(11,2);not null"`
	WasKickbacked      bool      `json:"was_kickbacked"`
	Product            Product   `json:"product"`
}

// Describes information we get from email
type PurchaseData struct {
	ProductId     string
	ProductName   string
	ProductURL    string
	PurchasePrice float32
	PurchaseAt    time.Time
	SellerName    string
}

func (purchase *Purchase) getProduct() {
	var product Product
	DB.Where("id = ?", purchase.ProductId).First(&product)

	purchase.Product = product
}
