package mapper

import (
	"github.com/Josey34/goshop/delivery/http/dto/request"
	"github.com/Josey34/goshop/delivery/http/dto/response"
	"github.com/Josey34/goshop/domain/entity"
	ucproduct "github.com/Josey34/goshop/usecase/product"
)

func ToCreateProductInput(req request.CreateProductRequest) ucproduct.CreateProductInput {
	return ucproduct.CreateProductInput{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}
}

func ToProductResponse(p *entity.Product) response.ProductResponse {
	return response.ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price.Value(),
		Stock:       p.Stock,
		CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func ToProductListResponse(products []*entity.Product) []response.ProductResponse {
	result := make([]response.ProductResponse, len(products))
	for i, p := range products {
		result[i] = ToProductResponse(p)
	}
	return result
}
