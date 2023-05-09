package main

import (
	"log"

	"github.com/Astenna/Nubes/evaluation/hotel_baseline/models"
	"github.com/aws/aws-lambda-go/lambda"
)

func GetHotelsInCity(city string) ([]models.Hotel, error) {
	res, err := models.GetHotelsInCity(city)
	log.Printf("OUT LEN %d", len(res))
	return res, err
}

func main() {
	lambda.Start(GetHotelsInCity)
}
