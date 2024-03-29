package models

import (
	// External Packages
	"github.com/martini-contrib/binding"
	"net/http"
	"time"
	// Application Specific
	. "github.com/pfacheris/kickback/db"
)

type User struct {
	Id            int64      `json:"id"`
	Email         string     `json:"email" binding:"required" sql:"size:255;not null;unique"`
	LastMessageId string     `json:"last_message_id" sql:"size:255"`
	RefreshToken  string     `json:"-" binding:"required" sql:"size:255"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     time.Time  `json:"-"`
	Purchases     []Purchase `json:"purchases"`
}

// This method implements binding.Validator and is executed by the binding.Validate middleware
func (user User) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	if len(user.Email) < 3 {
		errors = append(errors, binding.Error{
			FieldNames: []string{"email"},
			Message:    "Too short; minimum 3 characters",
		})
	} else if len(user.Email) > 120 {
		errors = append(errors, binding.Error{
			FieldNames: []string{"email"},
			Message:    "Too long; maximum 255 characters",
		})
	}
	return errors
}

func (user *User) Get(id int64) error {
	if err := DB.First(user, id).Error; err != nil {
		return err
	}

	user.GetPurchases()

	return nil
}

func (user *User) GetPurchases() {
	purchases := []Purchase{}
	DB.Where("user_id = ?", user.Id).Find(&purchases)

	for i, p := range purchases {
		(&p).getProduct()
		purchases[i] = p
	}

	user.Purchases = purchases
}
