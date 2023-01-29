package main

import (
	"fmt"
	"net/http"
	"text/template"

	clib "github.com/Astenna/Nubes/movie_review/client_lib"
)

func categoryHandler(w http.ResponseWriter, r *http.Request) {
	categoryName := r.URL.Path[len("/category/"):]
	initializedCategory, err := clib.LoadCategory(categoryName)
	if initializedCategory == nil {
		fmt.Fprintf(w, "Category name %s not found", categoryName)
		return
	}
	if err != nil {
		fmt.Fprintf(w, "Error occurred when retrieving the category %s", err.Error())
		return
	}

	t, _ := template.ParseFiles("templates//category.html")
	t.Execute(w, initializedCategory)
}
