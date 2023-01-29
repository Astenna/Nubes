package main

import (
	"github.com/Astenna/Nubes/movie_review_baseline/faas/db"
	"github.com/Astenna/Nubes/movie_review_baseline/faas/models"
	"github.com/aws/aws-lambda-go/lambda"
)

func GetCategoryMovieListHandler(categoryName string) ([]models.CategoryListItem, error) {
	return db.GetMoviesByCategory(categoryName)
}

func main() {
	lambda.Start(GetCategoryMovieListHandler)
}
