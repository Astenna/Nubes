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

	return invokeLambda(reqBody, marshalledInput)
}

func invokeLambda(reqBody *GatewayParam, marshalledInput []byte) (interface{}, error) {
	switch reqBody.FunctionName {
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
	case "Delete":
		log.Printf("INVOKING: Delete")
		return LambdaClient.Invoke(&lambda.InvokeInput{
			FunctionName: aws.String("Delete"),
			Payload:      marshalledInput,
		})
	case "ReferenceGetStubs":
		log.Printf("INVOKING: ReferenceGetStubs (Get all hotels in a City OR Get Users' reservations)")
		return LambdaClient.Invoke(&lambda.InvokeInput{
			FunctionName: aws.String("ReferenceGetStubs"),
			Payload:      marshalledInput,
		})
	case "SetField":
		log.Printf("INVOKING: SetField (set hotel rate)")
		return LambdaClient.Invoke(&lambda.InvokeInput{
			FunctionName: aws.String("SetField"),
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
