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

func GatewayHandler(param map[string]interface{}) (interface{}, error) {
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
	case "CityGetAllHotels":
		log.Printf("INVOKING: CityGetAllHotels")
		return LambdaClient.Invoke(&lambda.InvokeInput{
			FunctionName: aws.String("CityGetAllHotels"),
			Payload:      marshalledInput,
		})

	case "UserVerifyPassword":
		log.Printf("INVOKING: UserVerifyPassword")
		return LambdaClient.Invoke(&lambda.InvokeInput{
			FunctionName: aws.String("UserVerifyPassword"),
			Payload:      marshalledInput,
		})

	case "CityGetHotelsCloseTo":
		log.Printf("INVOKING: CityGetHotelsCloseTo")
		return LambdaClient.Invoke(&lambda.InvokeInput{
			FunctionName: aws.String("CityGetHotelsCloseTo"),
			Payload:      marshalledInput,
		})

	case "CityGetHotelsWithBestRates":
		log.Printf("INVOKING: CityGetHotelsWithBestRates")
		return LambdaClient.Invoke(&lambda.InvokeInput{
			FunctionName: aws.String("CityGetHotelsWithBestRates"),
			Payload:      marshalledInput,
		})
	case "Export":
		log.Printf("INVOKING: Export")
		return LambdaClient.Invoke(&lambda.InvokeInput{
			FunctionName: aws.String("Export"),
			Payload:      marshalledInput,
		})

	default:
		log.Printf("DEFAULT: " + reqBody.FunctionName)
		return "", fmt.Errorf(reqBody.FunctionName + " not supported")

	}
}

func main() {
	lambdaHandler.Start(GatewayHandler)
}
