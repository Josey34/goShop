package order

import (
	"context"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/valueobject"
	"github.com/Josey34/goshop/repository"
)

type UpdateOrderStatusInput struct {
	ID     string
	Status valueobject.OrderStatus
}

type UpdateOrderUseCase struct {
	orderRepo repository.OrderRepository
}

func NewUpdateOrderStatusUseCase(repo repository.OrderRepository) *UpdateOrderUseCase {
	return &UpdateOrderUseCase{orderRepo: repo}
}

func (uc *UpdateOrderUseCase) Execute(ctx context.Context, input UpdateOrderStatusInput) (*entity.Order, error) {
	orderFound, err := uc.orderRepo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if err := orderFound.TransitionTo(input.Status); err != nil {
		return nil, err
	}

	if err := uc.orderRepo.UpdateStatus(ctx, input.ID, input.Status); err != nil {
		return nil, err
	}

	return orderFound, nil
}
