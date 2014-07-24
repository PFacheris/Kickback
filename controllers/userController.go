package controllers

import (
	// External Packages
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"net/http"
	"strconv"

	// Application Specific Imports
	. "github.com/pfacheris/kickback/db"
	. "github.com/pfacheris/kickback/models"
)

type UserController struct{}

func (controller UserController) Create(user User, errs binding.Errors, res http.ResponseWriter, r render.Render, params martini.Params) {
	// Check for input validation errors
	if errs.Len() > 0 {
		HandleError("json", 422, errs[0], r)
		return
	}

	// Attempt Save to DB and Handle Result
	if err := DB.Create(&user).Error; err != nil {
		HandleError("json", http.StatusConflict, err, r)
		return
	}

	// Return Result
	res.Header().Set("Location", fmt.Sprintf("/users/%d", user.Id))
	r.JSON(http.StatusOK, user)
}

func (controller UserController) Read(res http.ResponseWriter, r render.Render, params martini.Params) {
	// Return JSON
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Parse Query Param
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		HandleError("json", 422, err, r)
		return
	}

	// Read From DB
	user := User{}
	if err = user.Get(id); err != nil {
		HandleError("json", 404, err, r)
		return
	}

	// Return Result
	r.JSON(http.StatusOK, user)
}
