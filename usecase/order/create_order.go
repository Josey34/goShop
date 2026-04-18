package order

import (
	"context"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/factory"
	"github.com/Josey34/goshop/pkg/idgen"
	"github.com/Josey34/goshop/repository"
)

type CreateOrderInput struct {
	CustomerID string
	Items      []factory.CreateOrderItemInput
}

type CreateOrderUseCase struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
	orderQueue  repository.OrderQueueRepository
	idGen       idgen.IDGenerator
}

func NewCreateOrderUseCase(orderRepo repository.OrderRepository, productRepo repository.ProductRepository, orderQueue repository.OrderQueueRepository, idGen idgen.IDGenerator) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		orderQueue:  orderQueue,
		idGen:       idGen,
	}
}

func (uc *CreateOrderUseCase) Execute(ctx context.Context, input CreateOrderInput) (*entity.Order, error) {
	var orderItems []factory.CreateOrderItemInput
	for _, item := range input.Items {
		product, err := uc.productRepo.FindByID(ctx, item.ProductID)
		if err != nil {
			return nil, err
		}

		if err := product.ReduceStock(item.Quantity); err != nil {
			return nil, err
		}

		if err := uc.productRepo.Update(ctx, product); err != nil {
			return nil, err
		}

		orderItems = append(orderItems, factory.CreateOrderItemInput{
			ProductID: product.ID,
			Name:      product.Name,
			Price:     product.Price,
			Quantity:  item.Quantity,
		})
	}

	order, err := factory.NewOrder(uc.idGen.Generate(), factory.CreateOrderInput{
		CustomerID: input.CustomerID,
		Items:      orderItems,
	})
	if err != nil {
		return nil, err
	}

	if err := uc.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	if err := uc.orderQueue.Enqueue(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}
