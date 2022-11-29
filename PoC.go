package main

import (
	f "github.com/Astenna/Thesis_PoC/faas/types"
)

func Handler() string {
	return "Hello World!"
}

func main() {
	//var t = beldilib.LambdaClient
	//lambda.Start(Handler)
	var test = new(f.Shop)
	var user = test.Owner.Get()
	var _ = user.Id
}
