package models

import (
// External Packages
)

type Currency struct {
	Id     int64  `json:"id"`
	Name   string `json:"name" binding:"required" sql:"type:size:255;not null"`
	Symbol string `json:"symbol" binding:"required" sql:"type:size:1;not null;unique"`
}
