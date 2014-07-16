package userController

import (
  // External Packages
  "net/http"
  "encoding/json"
  "github.com/codegangsta/martini"
  "github.com/martini-contrib/render"

  // Application Specific Imports
  . "github.com/pfacheris/kickback/models/user"
)

func create(user User, res http.ResponseWriter, params martini.Params) {
  // Return JSON
  res.Header().Set("Content-Type", "application/json; charset=UTF-8")

  // Attempt Save to DB and Handle Result
  createdUser, err := user.Save()
  switch {
	case err != nil:
		return http.StatusConflict, json.Marshal(err)
	case err == nil:
		res.Header().Set("Location", fmt.Sprintf("/users/%d", id))
		return http.StatusCreated, json.Marshal(user)
	}
}


func read() {

}

func update() {

}

func destroy() {

}
