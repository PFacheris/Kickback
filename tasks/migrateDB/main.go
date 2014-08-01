package main

import (
	. "github.com/pfacheris/kickback/db"
	. "github.com/pfacheris/kickback/models"
)

func main() {
	DB.AutoMigrate(User{})
	DB.AutoMigrate(Purchase{})
	DB.AutoMigrate(Product{})
}
