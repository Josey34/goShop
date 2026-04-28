package main

import (
	"context"
	"math/rand"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(ctx context.Context, event ProcessPaymentEvent) (ProcessPaymentResponse, error) {
	if rand.Intn(10) < 8 {
		return ProcessPaymentResponse{
			OrderID:       event.OrderID,
			PaymentStatus: "success",
			Message:       "payment processed successfully",
		}, nil
	}

	return ProcessPaymentResponse{
		OrderID:       event.OrderID,
		PaymentStatus: "failed",
		Message:       "payment declined",
	}, nil

}
