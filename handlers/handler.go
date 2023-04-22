package handlers

import (
	"build-framework/framework"
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
	ctx.WriteString("list")
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

func PostsPageHandler(ctx *framework.HttpContext) {
	ctx.WriteString("form")
}

func PostsHandler(ctx *framework.HttpContext) {
	ctx.WriteString("post")
}
