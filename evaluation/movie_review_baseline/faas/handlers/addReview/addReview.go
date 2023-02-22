package main

import (
	"fmt"

	"github.com/Astenna/Nubes/evaluation/movie_review_baseline/faas/db"
	"github.com/Astenna/Nubes/evaluation/movie_review_baseline/faas/models"
	"github.com/aws/aws-lambda-go/lambda"
)

func AddReviewHandler(review models.Review) (string, error) {
	if review.MovieId == "" || review.ReviewerId == "" {
		return "", fmt.Errorf("MovieId and ReviewerId must be set")
	}
	result, err := db.Insert(review, "Review")
	if err != nil {
		return "", err
	}
	return result, nil
}

func main() {
	lambda.Start(AddReviewHandler)
}
