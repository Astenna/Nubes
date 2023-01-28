package main

import (
	"fmt"

	org "github.com/Astenna/Nubes/evaluation/movie_review/faas/types"

	lib "github.com/Astenna/Nubes/lib"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/mitchellh/mapstructure"
)

func ExportHandler(input aws.JSONValue) (string, error) {
	if input["TypeName"] == "" {
		return "", fmt.Errorf("missing TypeName in HandlerParameters")
	}

	switch input["TypeName"] {
	
	 case "Account":
			newAccount := new(org.Account)
			mapstructure.Decode(input["Parameter"], newAccount)
			return lib.Insert(newAccount)
	 
	 case "Category":
			newCategory := new(org.Category)
			mapstructure.Decode(input["Parameter"], newCategory)
			return lib.Insert(newCategory)
	 
	 case "Movie":
			newMovie := new(org.Movie)
			mapstructure.Decode(input["Parameter"], newMovie)
			return lib.Insert(newMovie)
	 
	 case "Review":
			newReview := new(org.Review)
			mapstructure.Decode(input["Parameter"], newReview)
			return lib.Insert(newReview)
	 

	default:
		return "", fmt.Errorf("%s not supported",  input["TypeName"])

	}
}

func main() {
	lambda.Start(ExportHandler)
}
