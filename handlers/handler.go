package handlers

import (
	"build-framework/framework"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
)

type StudentResponse struct {
	Name string `json:"name"`
}

func StudentsHandler(ctx *framework.HttpContext) {
	// query stringのnameを取り出す
	name := ctx.QueryKey("name", "")

	response := StudentResponse{
		Name: name,
	}

	ctx.Json(response)
}

func ListHandler(ctx *framework.HttpContext) {
	list := make([]string, 0)
	ctx.WriteString(list[0])
}

func ListItemHandler(ctx *framework.HttpContext) {
	ctx.WriteString("list item")
}

func ListItemPictureHandler(ctx *framework.HttpContext) {
	listID := ctx.GetParams(":list_id", "")
	pictureID := ctx.GetParams(":picture_id", "")

	type Output struct {
		ListID    string `json:"list_id"`
		PictureID string `json:"picture_id"`
	}

	ctx.Json(&Output{
		ListID:    listID,
		PictureID: pictureID,
	})
}

func UsersHandler(ctx *framework.HttpContext) {
	ctx.WriteString("users")
}

type PostForm struct {
	Name string
}

func FormHandler(ctx *framework.HttpContext) {
	authUser := ctx.Get("AuthUser", "defaultValue")
	postForm := &PostForm{
		Name: authUser.(string),
	}
	ctx.RenderHtml("./html/posts.html", postForm)
}

func PostsHandler(ctx *framework.HttpContext) {
	name := ctx.FormKey("name", "defaultName")
	age := ctx.FormKey("age", "20")
	fileInfo, err := ctx.FormFile("file")
	if err != nil {
		ctx.W.WriteHeader(http.StatusBadRequest)
	}

	ioutil.WriteFile(fmt.Sprintf("%s_%s_%s", name, age, fileInfo.Filename), fileInfo.Data, fs.ModePerm)
	if err != nil {
		ctx.W.WriteHeader(http.StatusInternalServerError)
	}

	ctx.WriteString("post")
}

type UserPost struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func UserPostHandler(ctx *framework.HttpContext) {
	userPost := &UserPost{}
	if err := ctx.BindJson(userPost); err != nil {
		ctx.W.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.Json(userPost)
}

func FetchApiHandler(ctx *framework.HttpContext) {
	queryKey := ctx.QueryKey("callback", "cb")
	ctx.JsonP(queryKey, "")

}
