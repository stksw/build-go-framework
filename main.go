package main

import (
	"build-framework/framework"
	"build-framework/handlers"
)

func main() {
	e := framework.NewEngine()

	e.Router.Get("/list", handlers.ListHandler)
	e.Router.Get("/list/:id", handlers.ListItemHandler)
	e.Router.Get("/list/:list_id/picture/:picture_id", handlers.ListItemPictureHandler)
	e.Router.Get("/users", handlers.UsersHandler)
	e.Router.Get("/students", handlers.StudentsHandler)

	e.Router.Get("/form", handlers.FormHandler)
	e.Router.Post("/posts", handlers.PostsHandler)

	e.Router.Get("/fetch_api", handlers.FetchApiHandler)
	e.Run()
}
