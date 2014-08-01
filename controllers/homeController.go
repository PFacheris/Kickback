package controllers

import (
	"fmt"
	// External Packages
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/render"
	"net/http"
	// "code.google.com/p/google-api-go-client/gmail/v1"

	// Application Specific Imports
	. "github.com/pfacheris/kickback/db"
	. "github.com/pfacheris/kickback/models"
)

type HomeController struct{}

type userInfo struct {
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
}

func (controller HomeController) Dashboard(tokens oauth2.Tokens, r render.Render) {
	// Check if the user already exists
	email, err := getCurrentUserEmail(tokens.Access())
	if err != nil {
		HandleError("html", 500, err, r)
		return
	}

	user := User{}
	if err = DB.Where("email = ?", email).First(&user).Error; err != nil {
		// Check for err type
		if err != gorm.RecordNotFound {
			HandleError("html", 500, err, r)
			return
		}
		// User did not previously exist, create it
		fmt.Println(tokens.Refresh())
		fmt.Println(tokens.Access())
		user = User{
			Email:        email,
			RefreshToken: tokens.Refresh(),
		}

		if err = DB.Create(&user).Error; err != nil {
			HandleError("html", 500, err, r)
			return
		}

		// User created, render success page
		r.HTML(200, "dashboard", nil)
		return
	}
	// User previously existed, render success page
	r.HTML(200, "dashboard", nil)
}

func (controller HomeController) Landing(tokens oauth2.Tokens, r render.Render) {
	if !tokens.IsExpired() {
		// User is logged in, redirect to dashboard.
		r.Redirect("/dashboard")
		return
	}

	r.HTML(200, "landing", nil)
}

// Utility Functions
func getCurrentUserEmail(accessToken string) (string, error) {
	url := "https://www.googleapis.com/oauth2/v2/userinfo"
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	var currentUserInfo userInfo
	err = json.NewDecoder(res.Body).Decode(&currentUserInfo)
	if err != nil {
		return "", err
	}

	return currentUserInfo.Email, nil
}
