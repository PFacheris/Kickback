package main

import (

	// App-level
	. "github.com/pfacheris/kickback/db"
	"github.com/pfacheris/kickback/models"
	"github.com/pfacheris/kickback/tasks/lib/updatePurchaseForUser"
)

func main() {

	users := []models.User{}
	DB.Where("refresh_token != ''").Find(&users)

	for _, user := range users {
		updatePurchaseForUser.Do(&user)
	}

}
