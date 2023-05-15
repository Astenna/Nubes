package main

import (
	"github.com/Astenna/Nubes/evaluation/hotel_baseline/models"
	"github.com/aws/aws-lambda-go/lambda"
)

func ReserveRoom(param models.ReserveParam) error {
	return models.ReserveRoom(param)
}

func main() {
	lambda.Start(ReserveRoom)
}
