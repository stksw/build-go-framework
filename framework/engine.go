package framework

import (
	"build-framework/handlers"
	"net/http"
)

type Engine struct {
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if r.URL.Path == "/list" {
			handlers.ListHandler(w, r)
			return
		}

		if r.URL.Path == "/users" {
			handlers.UsersHandler(w, r)
			return
		}

		if r.URL.Path == "/students" {
			handlers.StudentHandler(w, r)
			return
		}
	}
}

func (e *Engine) Run() {
	http.ListenAndServe("localhost:8000", e)
}
