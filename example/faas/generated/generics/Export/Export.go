package main

import (
	"fmt"

	org "github.com/Astenna/Nubes/example/faas/types"

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
	
	 case "Discount":
			newDiscount := new(org.Discount)
			mapstructure.Decode(input["Parameter"], newDiscount)
			return lib.Insert(newDiscount)
	 
	 case "Order":
			newOrder := new(org.Order)
			mapstructure.Decode(input["Parameter"], newOrder)
			return lib.Insert(newOrder)
	 
	 case "Product":
			newProduct := new(org.Product)
			mapstructure.Decode(input["Parameter"], newProduct)
			return lib.Insert(newProduct)
	 
	 case "Shipping":
			newShipping := new(org.Shipping)
			mapstructure.Decode(input["Parameter"], newShipping)
			return lib.Insert(newShipping)
	 
	 case "Shop":
			newShop := new(org.Shop)
			mapstructure.Decode(input["Parameter"], newShop)
			return lib.Insert(newShop)
	 
	 case "User":
			newUser := new(org.User)
			mapstructure.Decode(input["Parameter"], newUser)
			return lib.Insert(newUser)
	 

	default:
		return "", fmt.Errorf("%s not supported",  input["TypeName"])

	}
}

func main() {
	lambda.Start(ExportHandler)
}
