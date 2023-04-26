package main

import (
	"github.com/Astenna/Nubes/evaluation/hotel_baseline/models"
	"github.com/aws/aws-lambda-go/lambda"
)

type RecommendHotelsParam struct {
	Count int
	City  string
}

func RecommendHotelsLocation(params RecommendHotelsParam) ([]models.Hotel, error) {
	return models.RecommendHotelsPrice(params.City, params.Count)
}

func main() {
	lambda.Start(RecommendHotelsLocation)
}
