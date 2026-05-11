package sqlite_test

import (
	"context"
	"testing"
	"time"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/errors"
	"github.com/Josey34/goshop/domain/valueobject"
	"github.com/Josey34/goshop/repository/sqlite"
)

func makeProduct(id, name string, price int64, stock int) *entity.Product {
	now := time.Now()
	return &entity.Product{
		ID:        id,
		Name:      name,
		Price:     valueobject.NewMoney(price),
		Stock:     stock,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func TestProductRepo_CreateAndFindByID(t *testing.T) {
	db := newTestDB(t)
	repo := sqlite.NewProductRepo(db)
	ctx := context.Background()

	p := makeProduct("p1", "Shirt", 1000, 10)
	if err := repo.Create(ctx, p); err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := repo.FindByID(ctx, "p1")
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if got.Name != "Shirt" {
		t.Errorf("got name=%s, want=Shirt", got.Name)
	}
	if got.Price.Value() != 1000 {
		t.Errorf("got price=%d, want=1000", got.Price.Value())
	}
	if got.Stock != 10 {
		t.Errorf("got stock=%d, want=10", got.Stock)
	}
}

func TestProductRepo_FindByID_NotFound(t *testing.T) {
	db := newTestDB(t)
	repo := sqlite.NewProductRepo(db)

	_, err := repo.FindByID(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.IsNotFound(err) {
		t.Errorf("expected NotFound error, got %v", err)
	}
}

func TestProductRepo_FindAll_Pagination(t *testing.T) {
	db := newTestDB(t)
	repo := sqlite.NewProductRepo(db)
	ctx := context.Background()

	for i, name := range []string{"A", "B", "C"} {
		repo.Create(ctx, makeProduct("p"+string(rune('1'+i)), name, 1000, 5))
	}

	products, err := repo.FindAll(ctx, valueobject.NewPagination(1, 2))
	if err != nil {
		t.Fatalf("findAll: %v", err)
	}
	if len(products) != 2 {
		t.Errorf("got %d, want=2", len(products))
	}
}

func TestProductRepo_Update(t *testing.T) {
	db := newTestDB(t)
	repo := sqlite.NewProductRepo(db)
	ctx := context.Background()

	p := makeProduct("p1", "Old", 500, 5)
	repo.Create(ctx, p)

	p.Name = "New"
	p.Price = valueobject.NewMoney(2000)
	p.UpdatedAt = time.Now()
	if err := repo.Update(ctx, p); err != nil {
		t.Fatalf("update: %v", err)
	}

	got, _ := repo.FindByID(ctx, "p1")
	if got.Name != "New" {
		t.Errorf("got name=%s, want=New", got.Name)
	}
	if got.Price.Value() != 2000 {
		t.Errorf("got price=%d, want=2000", got.Price.Value())
	}
}

func TestProductRepo_Delete(t *testing.T) {
	db := newTestDB(t)
	repo := sqlite.NewProductRepo(db)
	ctx := context.Background()

	repo.Create(ctx, makeProduct("p1", "Shirt", 1000, 5))
	if err := repo.Delete(ctx, "p1"); err != nil {
		t.Fatalf("delete: %v", err)
	}

	_, err := repo.FindByID(ctx, "p1")
	if err == nil {
		t.Error("expected error after delete")
	}
}
