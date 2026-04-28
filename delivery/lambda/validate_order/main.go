package main

import (
	"github.com/Josey34/goshop/repository/memory"
	"github.com/Josey34/goshop/usecase/order"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	repo := memory.NewOrderRepo()
	uc := order.NewGetOrderUseCase(repo)
	h := NewHandler(uc)
	lambda.Start(h.Handle)
}
