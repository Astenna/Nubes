package main

import (
	"github.com/Astenna/Nubes/evaluation/hotel_baseline_simple/models"
	"github.com/aws/aws-lambda-go/lambda"
)

type RecommendHotelsRateParam struct {
	Count int
	City  string
}

func RecommendHotelsRate(params RecommendHotelsRateParam) ([]models.Hotel, error) {
	return models.RecommendHotelsRate(params.City, params.Count)
}

func main() {
	lambda.Start(RecommendHotelsRate)
}
