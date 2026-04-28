package main

import (
	"context"

	"github.com/Josey34/goshop/domain/valueobject"
	ucorder "github.com/Josey34/goshop/usecase/order"
)

type Handler struct {
	getOrderUC    *ucorder.GetOrderUseCase
	updateOrderUC *ucorder.UpdateOrderUseCase
}

func NewHandler(uc *ucorder.GetOrderUseCase, updateUC *ucorder.UpdateOrderUseCase) *Handler {
	return &Handler{getOrderUC: uc, updateOrderUC: updateUC}
}

func (h *Handler) Handle(ctx context.Context, event FulfillOrderEvent) (FullfillOrderResponse, error) {
	o, err := h.getOrderUC.Execute(ctx, event.OrderID)
	if err != nil {
		return FullfillOrderResponse{}, err
	}

	if err := o.TransitionTo(valueobject.OrderStatusDelivered); err != nil {
		return FullfillOrderResponse{}, err
	}

	updated, err := h.updateOrderUC.Execute(ctx, ucorder.UpdateOrderStatusInput{
		ID:     event.OrderID,
		Status: valueobject.OrderStatusDelivered,
	})
	if err != nil {
		return FullfillOrderResponse{}, err
	}

	return FullfillOrderResponse{
		OrderID: event.OrderID,
		Status:  string(updated.Status),
		Message: "order fulfilled",
	}, nil

}
