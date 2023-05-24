package main

import (
	"github.com/Astenna/Nubes/evaluation/hotel_baseline/models"
	"github.com/aws/aws-lambda-go/lambda"
)

func GetUserReservations(userEmail string) ([]models.Reservation, error) {
	return models.GetUserReservations(userEmail)
}

func main() {
	lambda.Start(GetUserReservations)
}
