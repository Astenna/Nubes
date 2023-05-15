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
	reqBody := new(GatewayParam)

	bodyString := param["body"].(string)
	err := json.Unmarshal([]byte(bodyString), &reqBody)
	if err != nil {
		log.Printf("error occurred while umarshalling request body")
	}

	marshalledInput, err := json.Marshal(reqBody.Input)
	if err != nil {
		log.Printf("error occurred while arshalling lambda input")
	}

	switch reqBody.FunctionName {

	case "getHotelsInCity":
		log.Printf("INVOKING: getHotelsInCity")
		return LambdaClient.Invoke(&lambda.InvokeInput{
			FunctionName: aws.String("getHotelsInCity"), Payload: marshalledInput,
		})

	case "login":
		log.Printf("INVOKING: login")
		return LambdaClient.Invoke(&lambda.InvokeInput{
			FunctionName: aws.String("login"), Payload: marshalledInput,
		})

	case "recommendHotelsLocation":
		log.Printf("INVOKING: recommendHotelsLocation")
		return LambdaClient.Invoke(&lambda.InvokeInput{
			FunctionName: aws.String("recommendHotelsLocation"), Payload: marshalledInput,
		})

	case "recommendHotelsRate":
		log.Printf("INVOKING: recommendHotelsRate")
		return LambdaClient.Invoke(&lambda.InvokeInput{
			FunctionName: aws.String("recommendHotelsRate"),
			Payload:      marshalledInput,
		})

	case "reserveRoom":
		log.Printf("INVOKING: reserveRoom")
		return LambdaClient.Invoke(&lambda.InvokeInput{
			FunctionName: aws.String("reserveRoom"),
			Payload:      marshalledInput,
		})

	default:
		log.Printf("DEFAULT: " + reqBody.FunctionName)
		return "", fmt.Errorf("%s not supported", reqBody.FunctionName)
	}
}

func main() {
	lambdaHandler.Start(GatewayHandler)
}
