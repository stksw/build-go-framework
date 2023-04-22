package framework

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HttpContext struct {
	w      http.ResponseWriter
	r      *http.Request
	params map[string]string
}

func NewContext(w http.ResponseWriter, r *http.Request) *HttpContext {
	return &HttpContext{
		w:      w,
		r:      r,
		params: map[string]string{},
	}
}

func (ctx *HttpContext) Json(data any) {
	response, err := json.Marshal(data)
	if err != nil {
		ctx.w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.w.Header().Set("Content-Type", "application/json")
	ctx.w.WriteHeader(http.StatusOK)
	ctx.w.Write(response)
}

func (ctx *HttpContext) WriteString(data string) {
	fmt.Fprint(ctx.w, data)
}

func (ctx *HttpContext) QueryAll() map[string][]string {
	return ctx.r.URL.Query()
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

func (ctx *HttpContext) SetParams(dict map[string]string) {
	ctx.params = dict
}

func (ctx *HttpContext) GetParams(key string, defaultValue string) string {
	params := ctx.params
	if v, ok := params[key]; ok {
		return v
	}
	return defaultValue
}
