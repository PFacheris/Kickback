package models

import (
	// External Packages
	"time"
)

type Purchase struct {
	Id             int64   `json:"id"`
	PurchasePrice  float32 `json:"purchase_price" binding:"required" sql:"type:decimal(11,2);not null"`
	KickbackAmount float32 `json:"kickback_amount" binding:"required" sql:"type:decimal(11,2);not null"`
	PurchaseAt     time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time `json:"-"`
	UserId         int64     `json:"-"`
	ProductId      string    `json:"-"`
}

// Describes information we get from email
type PurchaseData struct {
	ProductId     string
	ProductName   string
	ProductURL    string
	PurchasePrice float32
	PurchaseAt    time.Time
}
