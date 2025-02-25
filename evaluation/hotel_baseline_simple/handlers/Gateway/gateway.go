package main

import (
	"encoding/json"
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

	case "getHotelsInCity": // Simple override
		log.Printf("INVOKING: getHotelsInCitySimple")
		return LambdaClient.Invoke(&lambda.InvokeInput{
			FunctionName: aws.String("getHotelsInCitySimple"), Payload: marshalledInput,
		})

	case "recommendHotelsLocation": // Simple override
		log.Printf("INVOKING: recommendHotelsLocationSimple")
		return LambdaClient.Invoke(&lambda.InvokeInput{
			FunctionName: aws.String("recommendHotelsLocationSimple"), Payload: marshalledInput,
		})

	case "recommendHotelsRate": // Simple override
		log.Printf("INVOKING: recommendHotelsRateSimple")
		return LambdaClient.Invoke(&lambda.InvokeInput{
			FunctionName: aws.String("recommendHotelsRateSimple"),
			Payload:      marshalledInput,
		})

	// case "getHotelsInCity":
	// 	log.Printf("INVOKING: getHotelsInCity")
	// 	return LambdaClient.Invoke(&lambda.InvokeInput{
	// 		FunctionName: aws.String("getHotelsInCity"), Payload: marshalledInput,
	// 	})

	case "login":
		log.Printf("INVOKING: login")
		return LambdaClient.Invoke(&lambda.InvokeInput{
			FunctionName: aws.String("login"), Payload: marshalledInput,
		})

	// case "recommendHotelsLocation":
	// 	log.Printf("INVOKING: recommendHotelsLocation")
	// 	return LambdaClient.Invoke(&lambda.InvokeInput{
	// 		FunctionName: aws.String("recommendHotelsLocation"), Payload: marshalledInput,
	// 	})

	// case "recommendHotelsRate":
	// 	log.Printf("INVOKING: recommendHotelsRate")
	// 	return LambdaClient.Invoke(&lambda.InvokeInput{
	// 		FunctionName: aws.String("recommendHotelsRate"),
	// 		Payload:      marshalledInput,
	// 	})

	case "reserveRoom":
		log.Printf("INVOKING: reserveRoom")
		return LambdaClient.Invoke(&lambda.InvokeInput{
			FunctionName: aws.String("reserveRoom"),
			Payload:      marshalledInput,
		})

	case "deleteUser":
		log.Printf("INVOKING: deleteUser")
		return LambdaClient.Invoke(&lambda.InvokeInput{
			FunctionName: aws.String("deleteUser"),
			Payload:      marshalledInput,
		})

	case "registerUser":
		log.Printf("INVOKING: registerUser")
		return LambdaClient.Invoke(&lambda.InvokeInput{
			FunctionName: aws.String("registerUser"),
			Payload:      marshalledInput,
		})

	case "setHotelRate":
		log.Printf("INVOKING: setHotelRate")
		return LambdaClient.Invoke(&lambda.InvokeInput{
			FunctionName: aws.String("setHotelRate"),
			Payload:      marshalledInput,
		})

	case "getUserReservations":
		log.Printf("INVOKING: getUserReservations")
		return LambdaClient.Invoke(&lambda.InvokeInput{
			FunctionName: aws.String("getUserReservations"),
			Payload:      marshalledInput,
		})

	default:
		log.Printf("DEFAULT: " + reqBody.FunctionName)
		return LambdaClient.Invoke(&lambda.InvokeInput{
			FunctionName: aws.String("___AAA___"),
			Payload:      marshalledInput,
		})
	}
}

func main() {
	lambdaHandler.Start(GatewayHandler)
}
