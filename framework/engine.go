package framework

import (
	"fmt"
	"net/http"
	"strings"
)

type Engine struct {
	Router *Router
}

type Router struct {
	routingTables map[string]*TreeNode
	middlewares   []func(ctx *HttpContext)
	noRoute       func(ctx *HttpContext)
}

func NewEngine() *Engine {
	return &Engine{
		Router: &Router{
			routingTables: map[string]*TreeNode{
				"get":     Constructor(),
				"post":    Constructor(),
				"put":     Constructor(),
				"patch":   Constructor(),
				"delete":  Constructor(),
				"options": Constructor(),
			},
			middlewares: []func(ctx *HttpContext){},
		},
	}
}

func (r *Router) Use(middleware func(ctx *HttpContext)) {
	r.middlewares = append(r.middlewares, middleware)
}

func (r *Router) UseNoRoute(handler func(ctx *HttpContext)) {
	r.noRoute = handler
}

func (r *Router) register(method string, pathname string, handler func(ctx *HttpContext)) error {
	routingTable := r.routingTables[method]
	// 末端の/を除く
	pathname = strings.TrimSuffix(pathname, "/")
	exist := routingTable.Search(pathname)
	if exist != nil {
		panic("already exist!")
	}

	routingTable.Insert(pathname, handler)
	return nil
}

// http getを登録する処理
func (r *Router) Get(pathname string, handler func(ctx *HttpContext)) error {
	return r.register("get", pathname, handler)
}

// http postを登録する処理
func (r *Router) Post(pathname string, handler func(ctx *HttpContext)) error {
	return r.register("post", pathname, handler)
}

// http putを登録する処理
func (r *Router) Put(pathname string, handler func(ctx *HttpContext)) error {
	return r.register("put", pathname, handler)
}

// http patchを登録する処理
func (r *Router) Patch(pathname string, handler func(ctx *HttpContext)) error {
	return r.register("patch", pathname, handler)
}

// http deleteを登録する処理
func (r *Router) Delete(pathname string, handler func(ctx *HttpContext)) error {
	return r.register("delete", pathname, handler)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := NewContext(w, r)
	ctx.Set("AuthUser", "test")

	routingTable := e.Router.routingTables[strings.ToLower(r.Method)]
	pathname := r.URL.Path
	pathname = strings.TrimSuffix(pathname, "/")
	targetNode := routingTable.Search(pathname)

	var targetHandler func(ctx *HttpContext)
	if targetNode == nil || targetNode.handler == nil {
		// 404 not found
		fmt.Println("not found")
		targetHandler = e.Router.noRoute
		if targetHandler == nil {
			// 事前にnot found handlerが登録されていなければ
			targetHandler = DefaultNotFoundHandler
		}
	} else {
		targetHandler = targetNode.handler
		// :list_id, :picture_idなどのparamsをdict型で返す
		paramDict := targetNode.ParseParams(r.URL.Path)
		// 引数をctxのparamsにセットする
		ctx.SetParams(paramDict)
	}

	handlers := append(e.Router.middlewares, targetHandler)
	ctx.SetHandler(handlers)
	ctx.Next()
}

func (e *Engine) Run() {
	http.ListenAndServe("localhost:8000", e)
}

func DefaultNotFoundHandler(ctx *HttpContext) {
	ctx.ResponseWriter().WriteHeader(http.StatusNotFound)
}
