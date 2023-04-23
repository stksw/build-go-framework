package framework

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/context"
)

type Engine struct {
	Router *Router
}

type Router struct {
	routingTables map[string]*TreeNode
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
		},
	}
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
	if targetNode == nil || targetNode.handler == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// :list_id, :picture_idなどのparamsをdict型で返す
	paramDict := targetNode.ParseParams(r.URL.Path)

	// 引数をctxのparamsにセットする
	ctx.SetParams(paramDict)

	ch := make(chan struct{})
	go func() {
		// タイムアウトを試す際にはコメントアウトを外す
		// time.Sleep(time.Second * 10)
		targetNode.handler(ctx)
		ch <- struct{}{}
	}()

	// 5秒待機する処理
	durationContext, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	// 5秒待ってもgoroutineが応答なければ、Doneを実行
	select {
	case <-durationContext.Done():
		ctx.SetHasTimeout(true)
		fmt.Println("timeout")
		ctx.W.Write([]byte("timeout"))
	case <-ch:
		fmt.Println("finish")
	}
}

func (e *Engine) Run() {
	http.ListenAndServe("localhost:8000", e)
}
