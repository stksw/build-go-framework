package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type StudentResponse struct {
	Name string `json:"name"`
}

func StudentsHandler(w http.ResponseWriter, r *http.Request) {

	queries := r.URL.Query()
	// query stringのnameを取り出す
	name := queries.Get("name")

	response := StudentResponse{
		Name: name,
	}

	responseData, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}

func ListHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "lists")
}

func ListItemHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "item")
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "users")
}
