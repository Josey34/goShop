package main

import (
	"context"

	"github.com/Josey34/goshop/usecase/order"
)

type Handler struct {
	getOrderUC *order.GetOrderUseCase
}

func NewHandler(uc *order.GetOrderUseCase) *Handler {
	return &Handler{getOrderUC: uc}
}

func (h *Handler) Handle(ctx context.Context, event ValidateOrderEvent) (ValidateOrderResponse, error) {
	o, err := h.getOrderUC.Execute(ctx, event.OrderID)
	if err != nil {
		return ValidateOrderResponse{
			OrderID: event.OrderID,
			IsValid: false,
			Errors:  []string{err.Error()},
		}, nil
	}
	if len(o.Items) == 0 {
		return ValidateOrderResponse{OrderID: event.OrderID, IsValid: false, Errors: []string{"no items"}}, nil
	}
	return ValidateOrderResponse{OrderID: event.OrderID, IsValid: true, Errors: []string{}}, nil
}
