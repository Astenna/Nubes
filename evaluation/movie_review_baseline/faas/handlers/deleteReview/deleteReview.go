package main

import (
	"fmt"

	"github.com/Astenna/Nubes/evaluation/movie_review_baseline/faas/db"
	"github.com/aws/aws-lambda-go/lambda"
)

func DeleteReview(reviewId string) error {
	if reviewId == "" {
		return fmt.Errorf("missing reviewId")
	}
	return db.Delete(reviewId, "Review")
}

func main() {
	lambda.Start(DeleteReview)
}
