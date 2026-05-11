package sqlite_test

import (
	"context"
	"testing"
	"time"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/valueobject"
	"github.com/Josey34/goshop/repository/sqlite"
)

func insertCustomerRaw(t *testing.T, repo *sqlite.CustomerRepo, ctx context.Context) {
	t.Helper()
	repo.Create(ctx, makeCustomer("c1", "john@example.com"))
}

func makeOrder(id, customerID string, status valueobject.OrderStatus) *entity.Order {
	now := time.Now()
	return &entity.Order{
		ID:         id,
		CustomerID: customerID,
		Items: []entity.OrderItem{
			{
				ID:        id + "-item-1",
				OrderID:   id,
				ProductID: "p1",
				Name:      "Shirt",
				Price:     valueobject.NewMoney(1000),
				Quantity:  2,
			},
		},
		Total:     valueobject.NewMoney(2000),
		Status:    status,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func TestOrderRepo_CreateAndFindByID(t *testing.T) {
	db := newTestDB(t)
	customerRepo := sqlite.NewCustomerRepo(db)
	orderRepo := sqlite.NewOrderRepo(db)
	ctx := context.Background()

	insertCustomerRaw(t, customerRepo, ctx)

	o := makeOrder("o1", "c1", valueobject.OrderStatusPending)
	if err := orderRepo.Create(ctx, o); err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := orderRepo.FindByID(ctx, "o1")
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if got.ID != "o1" {
		t.Errorf("got id=%s, want=o1", got.ID)
	}
	if got.Total.Value() != 2000 {
		t.Errorf("got total=%d, want=2000", got.Total.Value())
	}
	if got.Status != valueobject.OrderStatusPending {
		t.Errorf("got status=%s, want=PENDING", got.Status)
	}
	if len(got.Items) != 1 {
		t.Errorf("got %d items, want=1", len(got.Items))
	}
}

func TestOrderRepo_UpdateStatus(t *testing.T) {
	db := newTestDB(t)
	customerRepo := sqlite.NewCustomerRepo(db)
	orderRepo := sqlite.NewOrderRepo(db)
	ctx := context.Background()

	insertCustomerRaw(t, customerRepo, ctx)
	orderRepo.Create(ctx, makeOrder("o1", "c1", valueobject.OrderStatusPending))

	if err := orderRepo.UpdateStatus(ctx, "o1", valueobject.OrderStatusConfirmed); err != nil {
		t.Fatalf("update status: %v", err)
	}

	got, _ := orderRepo.FindByID(ctx, "o1")
	if got.Status != valueobject.OrderStatusConfirmed {
		t.Errorf("got status=%s, want=CONFIRMED", got.Status)
	}
}

func TestOrderRepo_FindByCustomer(t *testing.T) {
	db := newTestDB(t)
	customerRepo := sqlite.NewCustomerRepo(db)
	orderRepo := sqlite.NewOrderRepo(db)
	ctx := context.Background()

	insertCustomerRaw(t, customerRepo, ctx)
	orderRepo.Create(ctx, makeOrder("o1", "c1", valueobject.OrderStatusPending))
	orderRepo.Create(ctx, makeOrder("o2", "c1", valueobject.OrderStatusPending))

	orders, err := orderRepo.FindByCustomer(ctx, "c1", valueobject.NewPagination(1, 10))
	if err != nil {
		t.Fatalf("findByCustomer: %v", err)
	}
	if len(orders) != 2 {
		t.Errorf("got %d orders, want=2", len(orders))
	}
}
