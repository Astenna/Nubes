package main

import (
	"fmt"
	"github.com/Astenna/Thesis_PoC/faas/types"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"strconv"
)

func Handler4(shop types.Shop) (string, error) {
	if shop.Owner == nil {
		log.Println("owner is null!")
		fmt.Println("owner is null!")
		return "owner is null!", nil
	} else {
		log.Println(shop.Owner.Get().Id)
		fmt.Println(shop.Owner.Get().Id)
		return strconv.Itoa(shop.Owner.Get().Id), nil
	}
	return "", nil
	//if err := faas.CreateShop(shop); err != nil {
	//	//return err
	//}
	//return nil
}

func main() {
	lambda.Start(Handler4)
}
