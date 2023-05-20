package main

import (
	"log"

	"github.com/Astenna/Nubes/evaluation/hotel_baseline/models"
	"github.com/aws/aws-lambda-go/lambda"
)

func GetUserReservations(userEmail string) ([]models.Reservation, error) {
	log.Printf("Retrieving reservations of %s", userEmail)
	res, err := models.GetUserReservations(userEmail)
	log.Printf("%+v\n", res)
	if err != nil {
		log.Printf("%+v\n", err)
	}
	return res, err
}

func main() {
	lambda.Start(GetUserReservations)
}
