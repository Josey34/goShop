package main

import "github.com/aws/aws-lambda-go/lambda"

func main() {
	h := NewHandler()
	lambda.Start(h.Handle)
}
