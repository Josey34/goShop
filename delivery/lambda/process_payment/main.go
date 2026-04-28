package main

import "github.com/aws/aws-lambda-go/lambda"

func main() {
	h := &Handler{}
	lambda.Start(h.Handle)
}
