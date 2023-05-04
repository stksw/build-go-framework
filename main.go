package main

import (
	"build-framework/framework"
	"build-framework/handlers"
	"fmt"
	"net/http"
)

func main() {
	e := framework.NewEngine()
	router := e.Router

	router.Get("/list", handlers.ListHandler)
	router.Get("/list/:id", handlers.ListItemHandler)
	router.Get("/list/:list_id/picture/:picture_id", handlers.ListItemPictureHandler)
	router.Get("/users", handlers.UsersHandler)
	router.Get("/students", handlers.StudentsHandler)

	router.Get("/form", handlers.FormHandler)
	router.Post("/posts", handlers.PostsHandler)

	router.Get("/fetch_api", handlers.FetchApiHandler)

	router.Use(framework.AuthMiddleware)
	router.Use(framework.TimeCostMiddleware)
	router.Use(framework.TimeoutMiddleware)
	router.UseNoRoute(func(ctx *framework.HttpContext) {
		ctx.WriteString("not found....")
	})
	router.Use(framework.StaticFileMiddleware)

	e.Run()
}

type MyResponseWriter struct {
}

func (rw *MyResponseWriter) Header() http.Header {
	return nil
}

func (rw *MyResponseWriter) Write(data []byte) (int, error) {
	fmt.Println(string(data))
	return 0, nil
}

func (rw *MyResponseWriter) WriteHeader(statusCode int) {
	fmt.Println(statusCode)
}

func TimeCost(ctx *framework.HttpContext) {
	fmt.Println("timecost")
}

func AuthUser(ctx *framework.HttpContext) {
	fmt.Println("AuthUser")
}

func TimeOut(ctx *framework.HttpContext) {
	fmt.Println("timeout")
}

func Posts(ctx *framework.HttpContext) {
	ctx.WriteString("posts")
}
