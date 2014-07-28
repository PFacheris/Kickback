package main

import (
	"fmt"
	. "github.com/pfacheris/kickback/db"
	"github.com/pfacheris/kickback/models"
)

func main() {
	products := []models.Product{}
	DB.Find(&products)
	for _, product := range products {
		(&product).GetPurchases()
		fmt.Println(product.Name)
		fmt.Println(product.Purchases)
		fmt.Println()
	}
}
