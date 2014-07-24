package main

import (
	// External Packages
	"errors"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"

	// Application Specific Imports
	"github.com/pfacheris/kickback/controllers"
	. "github.com/pfacheris/kickback/models"
)

var m *martini.Martini

func main() {
	/*
	 * Load Martini Web Framework
	 */
	m = martini.New()

	// Setup middleware
	m.Use(martini.Static("static_dist"))
	m.Use(render.Renderer(render.Options{
		Directory:  "static_dist/views",
		Extensions: []string{".html"},
		IndentJSON: false,
	}))
	m.Use(martini.Recovery())
	m.Use(martini.Logger())
	m.Use(sessions.Sessions("my_session", sessions.NewCookieStore([]byte("secret123"))))
	m.Use(oauth2.Google(&oauth2.Options{
		ClientId:     "692789787338-t9f5805ou1uec14gl1l4fttohtkld54e.apps.googleusercontent.com",
		ClientSecret: "mfFPXn7ZZDlXw8TR2bwtgexD",
		RedirectURL:  "http://localhost:3000/oauth2callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/gmail.readonly", "https://www.googleapis.com/auth/gmail.compose"},
	}))

	// Setup routes
	r := martini.NewRouter()

	// Define Controller Instances
	homeController := controllers.HomeController{}
	userController := controllers.UserController{}

	r.Get("/", homeController.Index)

	r.Get("/users/:id", userController.Read)
	r.Post("/users", binding.Json(User{}), userController.Create)

	r.NotFound(func(r render.Render) {
		controllers.HandleError("html", 404, errors.New("Page not found."), r)
	})

	// Add the router action
	m.Action(r.Handle)
	m.Run()
}
