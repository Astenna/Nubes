package handlers

import (
	"fmt"
	"net/http"
	"text/template"

	clib "github.com/Astenna/Nubes/evaluation/movie_review/client_lib"
)

func CategoryHandler(w http.ResponseWriter, r *http.Request) {
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
	stubs, err := initializedCategory.Movies.GetStubs()
	if err != nil {
		fmt.Fprintf(w, "Error occurred when retrieving movies of the category %s", err.Error())
		return
	}

	templateInput := struct {
		Id     string
		Movies []clib.MovieStub
	}{Id: categoryName, Movies: stubs}
	t, _ := template.ParseFiles("templates//category.html")
	t.Execute(w, templateInput)
}
