package order_test

import (
	"context"
	"testing"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/factory"
	"github.com/Josey34/goshop/domain/valueobject"
	ucorder "github.com/Josey34/goshop/usecase/order"
)

func TestCreateOrderUseCase_Execute(t *testing.T) {
	t.Run("valid order reduces stock", func(t *testing.T) {
		productRepo := &mockProductRepo{
			products: map[string]*entity.Product{
				"p1": {ID: "p1", Name: "Shirt", Price: valueobject.NewMoney(1000), Stock: 10},
			},
		}
		uc := ucorder.NewCreateOrderUseCase(&mockOrderRepo{}, productRepo, &mockIDGen{id: "order-1"})

		order, err := uc.Execute(context.Background(), ucorder.CreateOrderInput{
			CustomerID: "c1",
			Items:      []factory.CreateOrderItemInput{{ProductID: "p1", Quantity: 2}},
		})

		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if order.ID != "order-1" {
			t.Errorf("got id=%s, want=order-1", order.ID)
		}
		if productRepo.products["p1"].Stock != 8 {
			t.Errorf("got stock=%d, want=8", productRepo.products["p1"].Stock)
		}
		if order.Total.Value() != 2000 {
			t.Errorf("got total=%d, want=2000", order.Total.Value())
		}
	})

	t.Run("insufficient stock returns error", func(t *testing.T) {
		productRepo := &mockProductRepo{
			products: map[string]*entity.Product{
				"p1": {ID: "p1", Name: "Shirt", Price: valueobject.NewMoney(1000), Stock: 1},
			},
		}
		uc := ucorder.NewCreateOrderUseCase(&mockOrderRepo{}, productRepo, &mockIDGen{id: "x"})

		_, err := uc.Execute(context.Background(), ucorder.CreateOrderInput{
			CustomerID: "c1",
			Items:      []factory.CreateOrderItemInput{{ProductID: "p1", Quantity: 5}},
		})
		if err == nil {
			t.Error("expected error for insufficient stock")
		}
	})
}
