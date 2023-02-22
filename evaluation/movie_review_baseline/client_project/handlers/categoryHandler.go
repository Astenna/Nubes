package handlers

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/Astenna/Nubes/evaluation/movie_review_baseline/client_project/models"
)

func CategoryHandler(w http.ResponseWriter, r *http.Request) {
	categoryName := r.URL.Path[len("/category/"):]

	list, err := invokeLambdaToGetList[models.CategoryListItem](categoryName, "getCategoryMovieList")
	if list == nil {
		fmt.Fprintf(w, "Category name %s not found", categoryName)
		return
	}
	if err != nil {
		fmt.Fprintf(w, "Error occurred when retrieving the category %s", err.Error())
		return
	}

	templInput := models.MoviesOfCategoryTemplateInput{Name: categoryName, Movies: list}
	t, _ := template.ParseFiles("templates//category.html")
	t.Execute(w, templInput)
}
