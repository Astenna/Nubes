package main

import (
	"github.com/Astenna/Nubes/movie_review_baseline/faas/db"
	"github.com/Astenna/Nubes/movie_review_baseline/faas/models"
	"github.com/aws/aws-lambda-go/lambda"
)

func GetMovieReviews(movieId string) ([]models.Review, error) {
	return db.GetMovieReviews(movieId)
}

func main() {
	lambda.Start(GetMovieReviews)
}
