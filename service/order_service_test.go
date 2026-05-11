package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/factory"
	"github.com/Josey34/goshop/domain/valueobject"
	"github.com/Josey34/goshop/repository/memory"
	"github.com/Josey34/goshop/service"
	ucorder "github.com/Josey34/goshop/usecase/order"
)

type mockQueue struct {
	enqueued []*entity.Order
	err      error
}

func (m *mockQueue) Enqueue(ctx context.Context, order *entity.Order) error {
	m.enqueued = append(m.enqueued, order)
	return m.err
}

type fixedIDGen struct{}

func (f *fixedIDGen) Generate() string { return "order-1" }

func buildOrderService(queue *mockQueue) *service.OrderService {
	orderRepo := memory.NewOrderRepo()
	productRepo := memory.NewProductRepo()
	idGen := &fixedIDGen{}

	createUC := ucorder.NewCreateOrderUseCase(orderRepo, productRepo, idGen)
	getUC := ucorder.NewGetOrderUseCase(orderRepo)
	listUC := ucorder.NewListProductUseCase(orderRepo)
	updateUC := ucorder.NewUpdateOrderStatusUseCase(orderRepo)

	return service.NewOrderService(createUC, getUC, listUC, updateUC, queue)
}

func seedProduct(ctx context.Context, repo interface {
	Create(context.Context, *entity.Product) error
}) {
	repo.Create(ctx, &entity.Product{
		ID:    "p1",
		Name:  "Shirt",
		Price: valueobject.NewMoney(1000),
		Stock: 10,
	})
}

func TestOrderService_CreateOrder_Enqueues(t *testing.T) {
	queue := &mockQueue{}
	svc := buildOrderService(queue)

	productRepo := memory.NewProductRepo()
	ctx := context.Background()
	seedProduct(ctx, productRepo)

	orderRepo := memory.NewOrderRepo()
	createUC := ucorder.NewCreateOrderUseCase(orderRepo, productRepo, &fixedIDGen{})
	getUC := ucorder.NewGetOrderUseCase(orderRepo)
	listUC := ucorder.NewListProductUseCase(orderRepo)
	updateUC := ucorder.NewUpdateOrderStatusUseCase(orderRepo)
	svc = service.NewOrderService(createUC, getUC, listUC, updateUC, queue)

	_, err := svc.CreateOrder(ctx, ucorder.CreateOrderInput{
		CustomerID: "c1",
		Items:      []factory.CreateOrderItemInput{{ProductID: "p1", Quantity: 1}},
	})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(queue.enqueued) != 1 {
		t.Errorf("got %d enqueued, want=1", len(queue.enqueued))
	}
}

func TestOrderService_CreateOrder_QueueFailureNonFatal(t *testing.T) {
	queue := &mockQueue{err: errors.New("sqs down")}

	productRepo := memory.NewProductRepo()
	ctx := context.Background()
	seedProduct(ctx, productRepo)

	orderRepo := memory.NewOrderRepo()
	createUC := ucorder.NewCreateOrderUseCase(orderRepo, productRepo, &fixedIDGen{})
	getUC := ucorder.NewGetOrderUseCase(orderRepo)
	listUC := ucorder.NewListProductUseCase(orderRepo)
	updateUC := ucorder.NewUpdateOrderStatusUseCase(orderRepo)
	svc := service.NewOrderService(createUC, getUC, listUC, updateUC, queue)

	order, err := svc.CreateOrder(ctx, ucorder.CreateOrderInput{
		CustomerID: "c1",
		Items:      []factory.CreateOrderItemInput{{ProductID: "p1", Quantity: 1}},
	})
	if err != nil {
		t.Fatalf("queue failure must be non-fatal, got err: %v", err)
	}
	if order == nil {
		t.Error("expected order returned even if queue fails")
	}
}
