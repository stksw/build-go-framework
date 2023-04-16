package framework

import (
	"net/http"
	"strings"
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
	// 末端の/を除く
	pathname = strings.TrimSuffix(pathname, "/")
	exist := r.routingTable.Search(pathname)
	if exist != nil {
		panic("already exist!")
	}

	r.routingTable.Insert(pathname, handler)
	return nil
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		pathname := r.URL.Path
		pathname = strings.TrimSuffix(pathname, "/")
		handler := e.Router.routingTable.Search(pathname)
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
