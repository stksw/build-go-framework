package main

import (
	"build-framework/framework"
	"build-framework/handlers"
)

func main() {
	e := framework.NewEngine()

	e.Router.Get("/list", handlers.ListHandler)
	e.Router.Get("/list/:id", handlers.ListItemHandler)
	// e.Router.Get("/list/name", handlers.ListItemHandler)
	e.Router.Get("/users", handlers.UsersHandler)
	e.Router.Get("/students", handlers.StudentsHandler)
	e.Run()
}
