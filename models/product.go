package models

import (
  // External Packages
  "time"
  // Application Specific
  . "github.com/pfacheris/kickback/db"
)

type Product struct {
  Id             string    `json:"id" binding:"required" sql:"type:size:255;not null;unique"`
  Name           string    `json:"name" binding:"required" sql:"type:size:255;not null"`
  URL            string    `json:"url" binding:"required" sql:"type:size:255;not null;unique"`
  CurrentPrice   float32   `json:"current_price" sql:"type:decimal(11,2)"`
  ScrapedAt      time.Time
  CreatedAt      time.Time
  UpdatedAt      time.Time
  DeletedAt      time.Time `json:"-"`
  Purchases      []Purchase
}

func (product *Product) Get(id int64) error {
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
