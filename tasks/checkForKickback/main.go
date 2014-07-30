package main

import (
	"github.com/kr/pretty"
	"github.com/pfacheris/kickback/config"
	. "github.com/pfacheris/kickback/db"
	"github.com/pfacheris/kickback/models"
	// "github.com/pfacheris/kickback/tasks/lib/kickbackemailer"
)

func main() {
	var users []models.User

	DB.Table("users").Select("users.id, users.email, users.refresh_token").Joins("join purchases on purchases.user_id=users.id").Where("purchases.was_kickbacked = ?", false).Where("purchases.kickback_amount > ?", config.KICKBACK_THRESHOLD).Scan(&users)

	// for every user that wasn't kickedback, and deserves one

	for _, user := range users {
		user.GetPurchases()

	}
}
