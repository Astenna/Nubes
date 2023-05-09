package main

import (
	"log"

	"github.com/Astenna/Nubes/evaluation/hotel_baseline/models"
	"github.com/aws/aws-lambda-go/lambda"
)

type RecommendHotelsParam struct {
	Count       int
	City        string
	Coordinates models.Coordinates
}

func RecommendHotelsLocation(params RecommendHotelsParam) ([]models.Hotel, error) {
	log.Printf("INPUT: %+v", params)
	res, err := models.RecommendHotelsLocation(params.City, params.Coordinates, params.Count)
	log.Printf("IOUT LEN %d", len(res))
	return res, err
}

func main() {
	lambda.Start(RecommendHotelsLocation)
}
