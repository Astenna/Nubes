package main

import (
	"fmt"
	//"github.com/aws/aws-lambda-go/lambda"
	"github.com/Astenna/Thesis_PoC/FaaSLib"
)

func Handler() string {
	return "Hello World!"
}

func main() {
	//var t = beldilib.LambdaClient
	//lambda.Start(Handler)
	var test = new(FaaSLib.Reference[string])
	fmt.Println(test.Id)
}
