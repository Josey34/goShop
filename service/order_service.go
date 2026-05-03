package service

import (
	"context"
	"log"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/valueobject"
	"github.com/Josey34/goshop/repository"
	ucorder "github.com/Josey34/goshop/usecase/order"
)

type OrderService struct {
	createUC       *ucorder.CreateOrderUseCase
	getUC          *ucorder.GetOrderUseCase
	listUC         *ucorder.ListOrderUseCase
	updateStatusUC *ucorder.UpdateOrderUseCase
	queue          repository.OrderQueueRepository
}

func NewOrderService(
	createUC *ucorder.CreateOrderUseCase,
	getUC *ucorder.GetOrderUseCase,
	listUC *ucorder.ListOrderUseCase,
	updateStatusUC *ucorder.UpdateOrderUseCase,
	queue repository.OrderQueueRepository,
) *OrderService {
	return &OrderService{createUC, getUC, listUC, updateStatusUC, queue}
}

func (s *OrderService) CreateOrder(ctx context.Context, input ucorder.CreateOrderInput) (*entity.Order, error) {
	order, err := s.createUC.Execute(ctx, input)
	if err != nil {
		return nil, err
	}

	for _, e := range order.PullEvents() {
		log.Printf("[event] %s", e.EventName())
	}

	if err := s.queue.Enqueue(ctx, order); err != nil {
		log.Printf("[event] failed to enqueue order: %v", err)
	}

	return order, nil
}

func (s *OrderService) GetByID(ctx context.Context, id string) (*entity.Order, error) {
	return s.getUC.Execute(ctx, id)
}

func (s *OrderService) ListByCustomer(ctx context.Context, customerID string, pagination valueobject.Pagination) ([]*entity.Order, error) {
	return s.listUC.Execute(ctx, customerID, pagination)
}

func (s *OrderService) UpdateStatus(ctx context.Context, input ucorder.UpdateOrderStatusInput) (*entity.Order, error) {
	return s.updateStatusUC.Execute(ctx, input)
}
