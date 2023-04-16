package framework

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HttpContext struct {
	Writer  http.ResponseWriter
	Request *http.Request
}

func NewContext(w http.ResponseWriter, r *http.Request) *HttpContext {
	return &HttpContext{
		Writer:  w,
		Request: r,
	}
}

func (ctx *HttpContext) Json(data any) {
	response, err := json.Marshal(data)
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.Writer.Header().Set("Content-Type", "application/json")
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Write(response)
}

func (ctx *HttpContext) WriteString(data string) {
	fmt.Fprint(ctx.Writer, data)
}

func (ctx *HttpContext) QueryAll() map[string][]string {
	return ctx.Request.URL.Query()
}

func (ctx *HttpContext) QueryKey(key string, defaultValue string) string {
	// クエリストリングをマップとして返す
	values := ctx.QueryAll()

	if target, ok := values[key]; ok {
		// mapのvalueが空だったら、defaultValueを返す
		if len(target) == 0 {
			return defaultValue
		}
		// mapのvalueが空でなければ、それを返す
		return target[len(target)-1]
	}

	// mapに該当するものがなければ、defaultValueを返す
	return defaultValue
}
