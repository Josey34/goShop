package main

import (
	"github.com/Josey34/goshop/repository/memory"
	ucorder "github.com/Josey34/goshop/usecase/order"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	repo := memory.NewOrderRepo()
	uc := ucorder.NewGetOrderUseCase(repo)
	h := NewHandler(uc)
	lambda.Start(h.Handle)
}
