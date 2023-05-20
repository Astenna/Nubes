package main

import (
	"github.com/Astenna/Nubes/evaluation/hotel_baseline/models"
	"github.com/aws/aws-lambda-go/lambda"
)

type SetHotelRateParam struct {
	Rate      float32
	CityName  string
	HotelName string
}

func SetHotelRate(param SetHotelRateParam) error {
	return models.SetHotelRate(param.CityName, param.HotelName, param.Rate)
}

func main() {
	lambda.Start(SetHotelRate)
}
