package main

import (
	"fmt"
	"net/http"
	"text/template"

	clib "github.com/Astenna/Nubes/movie_review/client_lib"
)

func movieHandler(w http.ResponseWriter, r *http.Request) {
	movieId := r.URL.Path[len("/movie/"):]
	initializedMovie, err := clib.LoadMovie(movieId)
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
