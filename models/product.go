package models

import (
	// External Packages
	"time"
	// Application Specific
	. "github.com/pfacheris/kickback/db"
)

type Product struct {
	Id            int64      `json:"id"`
	ProductId     string     `json:"productId" binding:"required" sql:"size:255;not null;unique"`
	Name          string     `json:"name" binding:"required" sql:"size:255;not null"`
	URL           string     `json:"url" binding:"required" sql:"size:255;not null;unique"`
	ScrapedAt     time.Time  `json:"scraped_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     time.Time  `json:"-"`
	Purchases     []Purchase `json:"purchases"`
	SmallImageURL string     `json:"small_image_url"`
}

func (product *Product) GetById(id int64) error {
	if err := DB.First(product, id).Error; err != nil {
		return err
	}

	purchases := []Purchase{}
	if err := DB.Model(product).Related(&purchases).Error; err != nil {
		return err
	}

	product.Purchases = purchases
	return nil
}

func (product *Product) Get(id string) error {
	if err := DB.Where("product_id = ?", id).First(product).Error; err != nil {
		return err
	}

	product.GetPurchases()
	return nil
}

func (product *Product) GetPurchases() {
	purchases := []Purchase{}
	DB.Where("product_id = ?", product.Id).Find(&purchases)

	product.Purchases = purchases
}
