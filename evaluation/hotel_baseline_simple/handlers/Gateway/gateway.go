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
		log.Printf("error occurred while marshalling lambda input")
	}

	switch reqBody.FunctionName {

	case "getHotelsInCitySimple":
		log.Printf("INVOKING: getHotelsInCitySimple")
		return LambdaClient.Invoke(&lambda.InvokeInput{
			FunctionName: aws.String("getHotelsInCitySimple"), Payload: marshalledInput,
		})

	case "recommendHotelsLocationSimple":
		log.Printf("INVOKING: recommendHotelsLocationSimple")
		return LambdaClient.Invoke(&lambda.InvokeInput{
			FunctionName: aws.String("recommendHotelsLocationSimple"), Payload: marshalledInput,
		})

	case "recommendHotelsRateSimple":
		log.Printf("INVOKING: recommendHotelsRateSimple")
		return LambdaClient.Invoke(&lambda.InvokeInput{
			FunctionName: aws.String("recommendHotelsRateSimple"),
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
