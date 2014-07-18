package models

import (
  // External Packages
)

type Domain struct {
  Id             int64     `json:"id"`
  URL            string    `json:"url" binding:"required" sql:"type:size:255;not null"`
  CurrencyId     int64     `json:"-"`
}
