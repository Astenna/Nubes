package main

import (
	"github.com/Astenna/Nubes/evaluation/hotel_baseline/models"
	"github.com/aws/aws-lambda-go/lambda"
)

func GetHotelsInCity(city string) ([]models.Hotel, error) {
	return models.GetHotelsInCity(city)
}

func main() {
	lambda.Start(GetHotelsInCity)
}
