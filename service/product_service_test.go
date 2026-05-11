package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/valueobject"
	"github.com/Josey34/goshop/pkg/idgen"
	"github.com/Josey34/goshop/repository/memory"
	"github.com/Josey34/goshop/service"
	ucproduct "github.com/Josey34/goshop/usecase/product"
)

func buildProductService() (*service.ProductService, *memory.ProductRepo) {
	repo := memory.NewProductRepo()
	fileStorage := memory.NewFileStorage()
	idGen := idgen.NewUUIDGenerator()

	createUC := ucproduct.NewCreateProductUseCase(repo, idGen)
	getUC := ucproduct.NewGetProductUseCase(repo)
	listUC := ucproduct.NewListProductUseCase(repo)
	updateUC := ucproduct.NewUpdateProductUseCase(repo)
	deleteUC := ucproduct.NewDeleteProductUseCase(repo)
	uploadUC := ucproduct.NewUploadProductImageUsecase(repo, fileStorage)

	return service.NewProductService(createUC, getUC, listUC, updateUC, deleteUC, uploadUC), repo
}

func seedProductInRepo(repo *memory.ProductRepo) {
	repo.Create(context.Background(), &entity.Product{
		ID:        "p1",
		Name:      "Shirt",
		Price:     valueobject.NewMoney(1000),
		Stock:     10,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
}

func TestProductService_Create(t *testing.T) {
	svc, _ := buildProductService()

	p, err := svc.Create(context.Background(), ucproduct.CreateProductInput{
		Name: "Shirt", Price: 1000, Stock: 10,
	})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if p.Name != "Shirt" {
		t.Errorf("got name=%s, want=Shirt", p.Name)
	}
}

func TestProductService_GetByID(t *testing.T) {
	svc, repo := buildProductService()
	seedProductInRepo(repo)

	p, err := svc.GetByID(context.Background(), "p1")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if p.ID != "p1" {
		t.Errorf("got id=%s, want=p1", p.ID)
	}
}

func TestProductService_List(t *testing.T) {
	svc, repo := buildProductService()
	seedProductInRepo(repo)

	products, err := svc.List(context.Background(), valueobject.NewPagination(1, 10))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(products) != 1 {
		t.Errorf("got %d products, want=1", len(products))
	}
}

func TestProductService_Update(t *testing.T) {
	svc, repo := buildProductService()
	seedProductInRepo(repo)

	p, err := svc.Update(context.Background(), ucproduct.UpdateProductInput{
		ID: "p1", Name: "Updated", Price: 2000, Stock: 5,
	})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if p.Name != "Updated" {
		t.Errorf("got name=%s, want=Updated", p.Name)
	}
}

func TestProductService_Delete(t *testing.T) {
	svc, repo := buildProductService()
	seedProductInRepo(repo)

	if err := svc.Delete(context.Background(), "p1"); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	_, err := svc.GetByID(context.Background(), "p1")
	if err == nil {
		t.Error("expected error after delete")
	}
}

func TestProductService_UploadImage(t *testing.T) {
	svc, repo := buildProductService()
	seedProductInRepo(repo)

	url, err := svc.UploadImage(context.Background(), ucproduct.UploadProductImage{
		ProductID:   "p1",
		Data:        []byte("imgdata"),
		ContentType: "image/jpeg",
	})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if url == "" {
		t.Error("expected non-empty url")
	}
}
