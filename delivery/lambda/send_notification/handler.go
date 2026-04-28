package main

import (
	"context"
	"log"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(ctx context.Context, event SendNotificationEvent) (SendNotificationResponse, error) {
	log.Printf("[notification] order=%s customer=%s msg=%s", event.OrderID, event.CustomerID, event.Message)

	return SendNotificationResponse{
		OrderID: event.OrderID,
		Sent:    true,
	}, nil
}
