package service

import (
	"context"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/valueobject"
	ucproduct "github.com/Josey34/goshop/usecase/product"
)

type ProductService struct {
	createUC *ucproduct.CreateProductUseCase
	getUC    *ucproduct.GetProductUseCase
	listUC   *ucproduct.ListProductsUseCase
	updateUC *ucproduct.UpdateProductUseCase
	deleteUC *ucproduct.DeleteProductUseCase
}

func NewProductService(
	createUC *ucproduct.CreateProductUseCase,
	getUC *ucproduct.GetProductUseCase,
	listUC *ucproduct.ListProductsUseCase,
	updateUC *ucproduct.UpdateProductUseCase,
	deleteUC *ucproduct.DeleteProductUseCase,
) *ProductService {
	return &ProductService{createUC, getUC, listUC, updateUC, deleteUC}
}

func (s *ProductService) Create(ctx context.Context, input ucproduct.CreateProductInput) (*entity.Product, error) {
	return s.createUC.Execute(ctx, input)
}

func (s *ProductService) GetByID(ctx context.Context, id string) (*entity.Product, error) {
	return s.getUC.Execute(ctx, id)
}

func (s *ProductService) List(ctx context.Context, pagination valueobject.Pagination) ([]*entity.Product, error) {
	return s.listUC.Execute(ctx, pagination)
}

func (s *ProductService) Update(ctx context.Context, input ucproduct.UpdateProductInput) (*entity.Product, error) {
	return s.updateUC.Execute(ctx, input)
}

func (s *ProductService) Delete(ctx context.Context, id string) error {
	return s.deleteUC.Execute(ctx, id)
}
