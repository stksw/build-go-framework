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

func UsersHandler(ctx *framework.HttpContext) {
	ctx.WriteString("users")
}
