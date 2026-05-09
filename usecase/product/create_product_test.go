package product_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/valueobject"
	ucproduct "github.com/Josey34/goshop/usecase/product"
)

// --- mocks ---

type mockProductRepo struct {
	createErr error
	created   *entity.Product
}

func (m *mockProductRepo) Create(ctx context.Context, p *entity.Product) error {
	m.created = p
	return m.createErr
}
func (m *mockProductRepo) FindByID(ctx context.Context, id string) (*entity.Product, error) {
	return nil, nil
}
func (m *mockProductRepo) FindAll(ctx context.Context, pagination valueobject.Pagination) ([]*entity.Product, error) {
	return nil, nil
}
func (m *mockProductRepo) Update(ctx context.Context, p *entity.Product) error { return nil }
func (m *mockProductRepo) Delete(ctx context.Context, id string) error         { return nil }

type mockIDGen struct{ id string }

func (m *mockIDGen) Generate() string { return m.id }

// --- tests ---

func TestCreateProductUseCase_Execute(t *testing.T) {
	cases := []struct {
		name      string
		input     ucproduct.CreateProductInput
		repoErr   error
		wantErr   bool
		wantName  string
	}{
		{
			name:     "valid product",
			input:    ucproduct.CreateProductInput{Name: "Shirt", Price: 1000, Stock: 10},
			wantName: "Shirt",
		},
		{
			name:    "missing name",
			input:   ucproduct.CreateProductInput{Name: "", Price: 1000},
			wantErr: true,
		},
		{
			name:    "repo error",
			input:   ucproduct.CreateProductInput{Name: "Shirt", Price: 1000},
			repoErr: errors.New("db error"),
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &mockProductRepo{createErr: tc.repoErr}
			idGen := &mockIDGen{id: "test-id"}
			uc := ucproduct.NewCreateProductUseCase(repo, idGen)

			product, err := uc.Execute(context.Background(), tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("got err=%v, wantErr=%v", err, tc.wantErr)
			}
			if !tc.wantErr {
				if product.Name != tc.wantName {
					t.Errorf("got name=%s, want=%s", product.Name, tc.wantName)
				}
				if product.ID != "test-id" {
					t.Errorf("got id=%s, want=test-id", product.ID)
				}
			}
		})
	}
}
