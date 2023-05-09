package main

import (
	"log"

	"github.com/Astenna/Nubes/evaluation/hotel_baseline/models"
	"github.com/aws/aws-lambda-go/lambda"
)

func ReserveRoom(param models.ReserveParam) error {
	log.Printf("INPUT: %+v", param)
	err := models.ReserveRoom(param)
	if err != nil {
		log.Printf("INPUT: %+v", err)
	} else {
		log.Printf("INPUT: SUCCESS")
	}
	return err
}

func main() {
	lambda.Start(ReserveRoom)
}
