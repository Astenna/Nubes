package main

import (
	"log"

	"github.com/Astenna/Nubes/evaluation/hotel_baseline/models"
	"github.com/aws/aws-lambda-go/lambda"
)

type RecommendHotelsRateParam struct {
	Count int
	City  string
}

func RecommendHotelsRate(params RecommendHotelsRateParam) ([]models.Hotel, error) {
	log.Printf("INPUT: %+v", params)
	res, err := models.RecommendHotelsRate(params.City, params.Count)
	log.Printf("IOUT LEN %d", len(res))
	return res, err
}

func main() {
	lambda.Start(RecommendHotelsRate)
}
