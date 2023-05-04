package main

import (
	"build-framework/framework"
	"build-framework/handlers"
	"context"
	"fmt"
	"time"
)

func main() {
	e := framework.NewEngine()

	e.Router.Get("/list", func(ctx *framework.HttpContext) {
		framework.TimeCostMiddleware(ctx, framework.AuthMiddleware(ctx, framework.TimeoutMiddleware(ctx, handlers.ListHandler)))(ctx)
	})
	e.Router.Get("/list/:id", func(ctx *framework.HttpContext) {
		successCh := make(chan struct{})
		panicCh := make(chan struct{})
		durationContext, cancel := context.WithTimeout(ctx.Request().Context(), time.Second*5)
		defer cancel()

		go func() {
			defer func() {
				if err := recover(); err != nil {
					panicCh <- struct{}{}
				}
			}()
			framework.AuthMiddleware(ctx, handlers.ListItemHandler)
			successCh <- struct{}{}
		}()

		select {
		case <-durationContext.Done():
			fmt.Println("success")
		case <-panicCh:
			ctx.WriteString("panic")
		case <-successCh:
			fmt.Println("success")
		}
	})
	e.Router.Get("/list/:list_id/picture/:picture_id", handlers.ListItemPictureHandler)
	e.Router.Get("/users", handlers.UsersHandler)
	e.Router.Get("/students", handlers.StudentsHandler)

	e.Router.Get("/form", handlers.FormHandler)
	e.Router.Post("/posts", handlers.PostsHandler)

	e.Router.Get("/fetch_api", handlers.FetchApiHandler)
	e.Run()
}
