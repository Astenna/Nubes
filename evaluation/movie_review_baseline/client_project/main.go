package main

import (
	"net/http"

	"github.com/Astenna/Nubes/evaluation/movie_review_baseline/client_project/handlers"
)

func main() {

	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/category/", handlers.CategoryHandler)
	http.HandleFunc("/movie/", handlers.MovieHandler)
	http.ListenAndServe(":8080", nil)
}
