package framework

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"net/textproto"
)

type HttpContext struct {
	W      http.ResponseWriter
	R      *http.Request
	params map[string]string
}

func NewContext(w http.ResponseWriter, r *http.Request) *HttpContext {
	return &HttpContext{
		W:      w,
		R:      r,
		params: map[string]string{},
	}
}

func (ctx *HttpContext) Json(data any) {
	response, err := json.Marshal(data)
	if err != nil {
		ctx.W.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.W.Header().Set("Content-Type", "application/json")
	ctx.W.WriteHeader(http.StatusOK)
	ctx.W.Write(response)
}

func (ctx *HttpContext) JsonP(callback string, parameter any) error {
	ctx.W.Header().Add("Content-Type", "application/javascript")
	callback = template.JSEscapeString(callback)

	_, err := ctx.W.Write([]byte(callback))
	if err != nil {
		return err
	}

	_, err = ctx.W.Write([]byte("("))
	if err != nil {
		return err
	}

	parameterData, err := json.Marshal(parameter)
	if err != nil {
		return err
	}

	_, err = ctx.W.Write(parameterData)
	if err != nil {
		return err
	}

	_, err = ctx.W.Write([]byte(")"))
	if err != nil {
		return err
	}

	return nil
}

func (ctx *HttpContext) WriteString(data string) {
	fmt.Fprint(ctx.W, data)
}

func (ctx *HttpContext) QueryAll() map[string][]string {
	return ctx.R.URL.Query()
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

func (ctx *HttpContext) FormKey(key string, defaultValue string) string {
	if ctx.R.Form == nil {
		ctx.R.ParseMultipartForm(32 << 20)
	}
	if vs := ctx.R.Form[key]; len(vs) > 0 {
		return vs[0]
	}
	return defaultValue
}

type FormFileInfo struct {
	Data     []byte
	Filename string
	Header   textproto.MIMEHeader
	Size     int64
}

func (ctx *HttpContext) FormFile(key string) (*FormFileInfo, error) {
	file, fileHeader, err := ctx.R.FormFile(key)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return &FormFileInfo{
		Data:     data,
		Filename: fileHeader.Filename,
		Header:   fileHeader.Header,
		Size:     fileHeader.Size,
	}, nil
}

func (ctx *HttpContext) BindJson(data any) error {
	byteData, err := io.ReadAll(ctx.R.Body)
	if err != nil {
		return err
	}

	// structをjsonに変換
	return json.Unmarshal([]byte{}, byteData)
}

func (ctx *HttpContext) RenderHtml(filepath string, data any) error {
	// templateを読み込む
	t, err := template.ParseFiles(filepath)
	if err != nil {
		return err
	}
	return t.Execute(ctx.W, data)
}
