package product

import (
	"context"

	"github.com/Josey34/goshop/repository"
)

type UploadProductImage struct {
	ProductID   string
	Data        []byte
	ContentType string
}

type UploadProductImageUsecase struct {
	productRepo repository.ProductRepository
	fileStorage repository.FileStorage
}

func NewUploadProductImageUsecase(repo repository.ProductRepository, fileStorage repository.FileStorage) *UploadProductImageUsecase {
	return &UploadProductImageUsecase{productRepo: repo, fileStorage: fileStorage}
}

func (uc *UploadProductImageUsecase) Execute(ctx context.Context, input UploadProductImage) (string, error) {
	product, err := uc.productRepo.FindByID(ctx, input.ProductID)
	if err != nil {
		return "", err
	}

	key := "products/" + product.ID + "/image.jpg"

	imageURL, err := uc.fileStorage.Upload(ctx, key, input.Data, input.ContentType)
	if err != nil {
		return "", err
	}

	product.ImageURL = imageURL

	err = uc.productRepo.Update(ctx, product)
	if err != nil {
		return "", err
	}

	return imageURL, nil
}
