package main

import (
	"github.com/Astenna/Nubes/evaluation/movie_review_baseline/faas/db"
	"github.com/Astenna/Nubes/evaluation/movie_review_baseline/faas/models"
	"github.com/aws/aws-lambda-go/lambda"
)

func GetMovieByIdHandler(movieId string) (models.Movie, error) {
	result, err := db.GetById[models.Movie](movieId, "Movie")
	if err != nil {
		return *new(models.Movie), err
	}
	return *result, nil
}

func main() {
	lambda.Start(GetMovieByIdHandler)
}
