package entity_test

import (
	"testing"
	"time"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/valueobject"
)

func TestCustomer_Fields(t *testing.T) {
	email, _ := valueobject.NewEmail("john@example.com")
	phone, _ := valueobject.NewPhone("08123456789")
	addr, _ := valueobject.NewAddress("Jl. A", "Jakarta", "DKI", "10220")

	c := &entity.Customer{
		ID:           "c1",
		Name:         "John",
		Email:        email,
		Phone:        phone,
		Address:      addr,
		PasswordHash: "hashed",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if c.ID != "c1" {
		t.Errorf("got id=%s, want=c1", c.ID)
	}
	if c.Email.Value() != "john@example.com" {
		t.Errorf("got email=%s", c.Email.Value())
	}
	if c.Phone.Value() != "08123456789" {
		t.Errorf("got phone=%s", c.Phone.Value())
	}
	if c.Address.City != "Jakarta" {
		t.Errorf("got city=%s", c.Address.City)
	}
}
