package controllers

import (
  // External Packages
  "fmt"
  "net/http"
  "encoding/json"
  "strconv"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/binding"

  // Application Specific Imports
  . "github.com/pfacheris/kickback/models"
  . "github.com/pfacheris/kickback/db"
)

type UserController struct {}

func (controller UserController) Create(user User, errs binding.Errors, res http.ResponseWriter, params martini.Params) (int, []byte) {
  var err error

  // Return JSON
  res.Header().Set("Content-Type", "application/json; charset=UTF-8")

  // Check for input validation errors
  if errs.Len() > 0 { return handleErrors(422, errs[0]) }

  // Attempt Save to DB and Handle Result
  err = DB.Create(&user).Error
  if err != nil { return handleErrors(http.StatusConflict, err) }

  // Return Result
  res.Header().Set("Location", fmt.Sprintf("/users/%d", user.Id))
  json, _ := json.Marshal(user)
  return http.StatusCreated, json
}


func (controller UserController) Read(res http.ResponseWriter, params martini.Params) (int, []byte) {
  var err error

  // Return JSON
  res.Header().Set("Content-Type", "application/json; charset=UTF-8")

  // Parse Query Param
  id, err := strconv.ParseInt(params["id"], 10, 64)
  if err != nil { return handleErrors(422, err) }

  // Read From DB
  user := User{}
  err = DB.First(&user, id).Error
  if err != nil { return handleErrors(404, err) }

  // Return Result
  json, _ := json.Marshal(user)
  return http.StatusOK, json
}

func (controller UserController) Update() {

}

func (controller UserController) Destroy() {

}

func handleErrors(status int, e error) (int, []byte) {
  json, _ := json.Marshal(map[string]string{
    "message": e.Error(),
  })
  return status, json
}
