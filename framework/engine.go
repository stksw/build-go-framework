package framework

import (
	"net/http"
)

type Engine struct {
	Router *Router
}

type Router struct {
	// routingTable map[string]func(w http.ResponseWriter, r *http.Request)
	routingTable TreeNode
}

func NewEngine() *Engine {
	return &Engine{
		Router: &Router{},
	}
}

// http getを登録する処理
func (r *Router) Get(pathname string, handler func(w http.ResponseWriter, r *http.Request)) error {
	r.routingTable.Insert(pathname, handler)
	return nil
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		path := r.URL.Path
		handler := e.Router.routingTable.Search(path)
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
