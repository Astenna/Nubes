package main

import (
	"github.com/Astenna/Nubes/evaluation/hotel_baseline/models"
	"github.com/aws/aws-lambda-go/lambda"
)

type RecommendHotelsParam struct {
	Count       int
	City        string
	Coordinates models.Coordinates
}

func RecommendHotelsLocation(params RecommendHotelsParam) ([]models.Hotel, error) {
	return models.RecommendHotelsLocation(params.City, params.Coordinates, params.Count)
}

func main() {
	lambda.Start(RecommendHotelsLocation)
}
