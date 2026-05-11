package sqlite_test

import (
	"context"
	"testing"
	"time"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/valueobject"
	"github.com/Josey34/goshop/repository/sqlite"
)

func makeCustomer(id, email string) *entity.Customer {
	e, _ := valueobject.NewEmail(email)
	p, _ := valueobject.NewPhone("08123456789")
	a, _ := valueobject.NewAddress("Jl. A", "Jakarta", "DKI", "10220")
	now := time.Now()
	return &entity.Customer{
		ID:           id,
		Name:         "John",
		Email:        e,
		Phone:        p,
		Address:      a,
		PasswordHash: "hashed",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func TestCustomerRepo_CreateAndFindByEmail(t *testing.T) {
	db := newTestDB(t)
	repo := sqlite.NewCustomerRepo(db)
	ctx := context.Background()

	c := makeCustomer("c1", "john@example.com")
	if err := repo.Create(ctx, c); err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := repo.FindByEmail(ctx, "john@example.com")
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if got.ID != "c1" {
		t.Errorf("got id=%s, want=c1", got.ID)
	}
	if got.Email.Value() != "john@example.com" {
		t.Errorf("got email=%s", got.Email.Value())
	}
}

func TestCustomerRepo_ExistsByEmail(t *testing.T) {
	db := newTestDB(t)
	repo := sqlite.NewCustomerRepo(db)
	ctx := context.Background()

	exists, _ := repo.ExistsByEmail(ctx, "john@example.com")
	if exists {
		t.Error("expected false before insert")
	}

	repo.Create(ctx, makeCustomer("c1", "john@example.com"))

	exists, _ = repo.ExistsByEmail(ctx, "john@example.com")
	if !exists {
		t.Error("expected true after insert")
	}
}

func TestCustomerRepo_DuplicateEmail(t *testing.T) {
	db := newTestDB(t)
	repo := sqlite.NewCustomerRepo(db)
	ctx := context.Background()

	repo.Create(ctx, makeCustomer("c1", "john@example.com"))
	err := repo.Create(ctx, makeCustomer("c2", "john@example.com"))
	if err == nil {
		t.Error("expected error for duplicate email")
	}
}

func TestCustomerRepo_FindByID(t *testing.T) {
	db := newTestDB(t)
	repo := sqlite.NewCustomerRepo(db)
	ctx := context.Background()

	repo.Create(ctx, makeCustomer("c1", "john@example.com"))

	got, err := repo.FindByID(ctx, "c1")
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if got.ID != "c1" {
		t.Errorf("got id=%s, want=c1", got.ID)
	}
}
