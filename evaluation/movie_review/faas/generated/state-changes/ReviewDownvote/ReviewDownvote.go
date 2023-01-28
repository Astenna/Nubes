package main

import (
	org "github.com/Astenna/Nubes/evaluation/movie_review/faas/types"
	"github.com/Astenna/Nubes/lib"
	"github.com/aws/aws-lambda-go/lambda"
)

func DownvoteHandler(input lib.HandlerParameters) (int, error) {
	instance := new(org.Review)
	instance.Id = input.Id
	instance.Init()

	result, _err := instance.Downvote(input.Parameter.(org.Account))
	if _err != nil {
		return result, _err
	}

	return result, _err
}

func main() {
	lambda.Start(DownvoteHandler)
}
