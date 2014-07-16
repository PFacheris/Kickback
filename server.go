package kickback

import (
  // External Packages
  "github.com/codegangsta/martini"
  "github.com/codegangsta/martini-contrib/binding"
  "github.com/martini-contrib/render"

  // Application Specific Imports
  . "github.com/pfacheris/kickback/models/user"
)

var m *martini.Martini
func init() {
  /*
   * Load Martini Web Framework
   */
  m = martini.New()

  // Setup middleware
  m.Use(martini.Recovery())
  m.Use(martini.Logger())
  m.Use(render.Renderer())
  m.Use(MapEncoder)

  // Setup routes
  r := martini.NewRouter()
  r.Get(`/users/:id`, binding.Bind(User{}), userController.read)
  r.Post(`/users`, binding.Bind(User{}), userController.create)
  r.Put(`/users/:id`, binding.Bind(User{}), userController.update)
  r.Delete(`/users/:id`, binding.Bind(User{}), userController.destroy)

  // Add the router action
  m.Action(r.Handle)
}
