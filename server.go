package main

import (
  // External Packages
  "github.com/go-martini/martini"
  "github.com/martini-contrib/binding"
  "github.com/martini-contrib/render"

  // Application Specific Imports
  . "github.com/pfacheris/kickback/models"
  "github.com/pfacheris/kickback/controllers"
)

var m *martini.Martini
func main() {
  /*
   * Load Martini Web Framework
   */
  m = martini.New()

  // Setup middleware
  m.Use(martini.Static("public"))
  m.Use(martini.Recovery())
  m.Use(martini.Logger())
  m.Use(render.Renderer())

  // Setup routes
  r := martini.NewRouter()
  r.Get(`/users/:id`, userController.Read)
  r.Post(`/users`, binding.Json(User{}), userController.Create)
  r.Put(`/users/:id`, binding.Json(User{}), userController.Update)
  r.Delete(`/users/:id`, binding.Json(User{}), userController.Destroy)

  // Add the router action
  m.Action(r.Handle)
  m.Run()
}
