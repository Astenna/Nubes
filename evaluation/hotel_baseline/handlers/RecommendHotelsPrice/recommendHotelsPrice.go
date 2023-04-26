package main

import (
	"github.com/Astenna/Nubes/evaluation/hotel_baseline/models"
	"github.com/aws/aws-lambda-go/lambda"
)

type RecommendHotelsPriceParam struct {
	Count int
	City  string
}

func RecommendHotelsPrice(params RecommendHotelsPriceParam) ([]models.Hotel, error) {
	return models.RecommendHotelsPrice(params.City, params.Count)
}

func main() {
	lambda.Start(RecommendHotelsPrice)
}
