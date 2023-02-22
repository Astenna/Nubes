package handlers

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/Astenna/Nubes/evaluation/movie_review_baseline/client_project/models"
)

func MovieHandler(w http.ResponseWriter, r *http.Request) {
	movieId := r.URL.Path[len("/movie/"):]
	initializedMovie, err := invokeLambdaToGetSingleItem[models.Movie](movieId, "getMovieById")
	if initializedMovie == nil {
		fmt.Fprintf(w, "Movie not found")
		return
	}
	if err != nil {
		fmt.Fprintf(w, "Error occurred when retrieving the category %s", err.Error())
		return
	}

	t, _ := template.ParseFiles("templates//movie.html")
	t.Execute(w, initializedMovie)
}
