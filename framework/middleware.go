package framework

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"golang.org/x/net/context"
)

func TimeoutMiddleware(ctx *HttpContext) {
	now := time.Now()
	successCh := make(chan struct{})
	panicCh := make(chan struct{})

	go func() {
		time.Sleep(time.Second * 1)
		defer func() {
			// panicが起きたら、関数が終了する前にpanicChに空のstructを送る
			if err := recover(); err != nil {
				panicCh <- struct{}{}
			}
		}()

		ctx.Next()
		successCh <- struct{}{}
	}()

	// 5秒待機する処理
	durationContext, cancel := context.WithTimeout(ctx.Request().Context(), time.Second*5)
	defer cancel()

	// 5秒待ってもgoroutineが応答なければ、Doneを実行
	select {
	case <-durationContext.Done():
		ctx.SetHasTimeout(true)
		fmt.Println("timeout")
		ctx.ResponseWriter().Write([]byte("timeout"))
	case <-successCh:
		fmt.Println(time.Since(now).Milliseconds())
		fmt.Println("finish")
	case <-panicCh:
		fmt.Println("panic")
		ctx.ResponseWriter().WriteHeader(500)
	}
}

func AuthMiddleware(ctx *HttpContext) {
	ctx.Set("Auth", "test")
}

func TimeCostMiddleware(ctx *HttpContext) {
	now := time.Now()
	ctx.Next()
	fmt.Println("time cost: ", time.Since(now).Milliseconds())
}

func StaticFileMiddleware(ctx *HttpContext) {
	fileServer := http.FileServer(http.Dir("./static"))
	pathname := ctx.Request().URL.Path
	pathname = strings.TrimSuffix(pathname, "/")
	fPath := path.Join("./static", pathname)
	fInfo, err := os.Stat(fPath)

	// ファイルが存在しているならサーバーを起動
	fExists := err == nil && !fInfo.IsDir()
	if fExists {
		fileServer.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
		ctx.Abort()
		return
	}

}
