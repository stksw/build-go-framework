package framework

import (
	"errors"
	"net/http"
)

type Engine struct {
	Router *Router
}

type Router struct {
	routingTable map[string]func(w http.ResponseWriter, r *http.Request)
}

func NewEngine() *Engine {
	return &Engine{
		Router: &Router{},
	}
}

func (r *Router) Get(pathname string, handler func(w http.ResponseWriter, r *http.Request)) error {
	if r.routingTable == nil {
		r.routingTable = make(map[string]func(w http.ResponseWriter, r *http.Request))
	}

	if r.routingTable[pathname] != nil {
		return errors.New("existed")
	}

	r.routingTable[pathname] = handler
	return nil
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		handler := e.Router.routingTable[r.URL.Path]
		if handler == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		handler(w, r)
		return
	}
}

func (e *Engine) Run() {
	http.ListenAndServe("localhost:8000", e)
}
