package main

import "github.com/Astenna/Nubes/evaluation/movie_review_baseline/faas/db"

func main() {
	db.GetMoviesByCategory("Drama")
}
