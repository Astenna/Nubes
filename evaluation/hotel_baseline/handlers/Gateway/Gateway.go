package main

import (
	"encoding/json"
	"fmt"
	"log"

	lambdaHandler "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))

var LambdaClient = lambda.New(sess)

type GatewayParam struct {
	FunctionName string
	Input        interface{}
}

func GatewayHandler(param aws.JSONValue) (interface{}, error) {
	reqBody := param["body"].(aws.JSONValue)
	eventJson, _ := json.Marshal(reqBody)
	log.Printf("INPUT: %s", eventJson)

	switch reqBody["FunctionName"] {

	case "GetHotelsInCity":
		return LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("CityGetAllHotels"), Payload: reqBody["Parameter"].([]byte)})

	case "Login":
		return LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("UserVerifyPassword"), Payload: reqBody["Parameter"].([]byte)})

	case "RecommendHotelsLocation":
		return LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("CityGetHotelsCloseTo"), Payload: reqBody["Parameter"].([]byte)})

	case "RecommendHotelsRate":
		return LambdaClient.Invoke(&lambda.InvokeInput{
			FunctionName: aws.String("CityGetHotelsWithBestRates"),
			Payload:      (reqBody["Parameter"].([]byte)),
		})

	default:
		return "", fmt.Errorf("%s not supported", reqBody["FunctionName"])

	}
}

func main() {
	lambdaHandler.Start(GatewayHandler)
}
