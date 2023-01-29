package main

import (
	"net/http"

	clib "github.com/Astenna/Nubes/movie_review/client_lib"
)

func main() {

	existingMovieId := "hand-written-uuid"
	newAccount := clib.AccountStub{Password: "password", Nickname: "mynick", Email: "test"}
	_, err := clib.ExportAccount(newAccount)
	_ = err
	newReview := clib.MovieReviewStub{Rating: 2, Text: "boring", Movie: clib.Reference[clib.MovieStub](existingMovieId), Reviewer: "test"}
	_, err = clib.ExportMovieReview(newReview)
	_ = err

	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/category/", categoryHandler)
	http.HandleFunc("/movie/", movieHandler)
	http.ListenAndServe(":8080", nil)
}
