package controllers

import (
  // External Packages
  "net/http"
  "encoding/json"
  "github.com/martini-contrib/oauth2"
  "github.com/martini-contrib/render"
  "github.com/jinzhu/gorm"
  // "code.google.com/p/google-api-go-client/gmail/v1"

  // Application Specific Imports
  . "github.com/pfacheris/kickback/models"
  . "github.com/pfacheris/kickback/db"
)

type HomeController struct {}

type userInfo struct {
  Email          string `json:"email"`
  VerifiedEmail  bool   `json:"verified_email"`
}

func (controller HomeController) Index(tokens oauth2.Tokens, r render.Render) {
  if tokens.IsExpired() {
    // User is not logged in
    // Render Landing Page HTML
    r.HTML(200, "landing", nil)
    return
  }

  // Check if the user already exists
  email, err := getCurrentUserEmail(tokens.Access())
  if err != nil {
    return
  }

  user := User{}
  if err = DB.Where("email = ?", email).First(&user).Error; err != nil {
    // Check for err type
    if err != gorm.RecordNotFound {
      return
    }
    // User did not previously exist, create it
    user = User{
      Email: email,
      AccessToken: tokens.Access(),
      RefreshToken: tokens.Refresh(),
      ExpireTokenAt: tokens.ExpiryTime(),
    }

    if err = DB.Create(&user).Error; err != nil {
      return
    }

    // User created, render success page
    r.HTML(200, "home", nil)
    return
  }

  // User previously existed, render success page
  r.HTML(200, "home", nil)
  return
}

// Utility Functions
func getCurrentUserEmail(accessToken string) (string, error) {
  url := "https://www.googleapis.com/oauth2/v2/userinfo"
  client := &http.Client{}
  req, _ := http.NewRequest("GET", url, nil)
  req.Header.Add("Authorization", "Bearer " + accessToken)

  res, _ := client.Do(req)

  var currentUserInfo userInfo
  err := json.NewDecoder(res.Body).Decode(&currentUserInfo)
  if err != nil {
    return "", err
  }

  return currentUserInfo.Email, nil
}

func handleHTMLErrors(status int, e error) (int, []byte) {
  json, _ := json.Marshal(map[string]string{
    "message": e.Error(),
  })
  return status, json
}
