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

func (h *Handler) Handle(ctx context.Context, event CalculateTotalEvent) (CalculateTotalResponse, error) {
	o, err := h.getOrderUC.Execute(ctx, event.OrderID)

	if err != nil {
		return CalculateTotalResponse{}, nil
	}

	subtotal := o.Total.Value()
	tax := o.Total.Percentage(11).Value()
	total := subtotal + tax

	return CalculateTotalResponse{
		OrderID:  event.OrderID,
		Subtotal: subtotal,
		Tax:      tax,
		Total:    total,
	}, nil
}
