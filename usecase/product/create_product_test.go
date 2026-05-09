package product_test

import (
	"context"
	"errors"
	"testing"

	ucproduct "github.com/Josey34/goshop/usecase/product"
)

func TestCreateProductUseCase_Execute(t *testing.T) {
	cases := []struct {
		name     string
		input    ucproduct.CreateProductInput
		repoErr  error
		wantErr  bool
		wantName string
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
			uc := ucproduct.NewCreateProductUseCase(repo, &mockIDGen{id: "test-id"})

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
