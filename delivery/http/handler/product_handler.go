package handler

import (
	"context"
	"io"
	"net/http"

	"github.com/Josey34/goshop/delivery/http/dto/mapper"
	"github.com/Josey34/goshop/delivery/http/dto/request"
	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/valueobject"
	ucproduct "github.com/Josey34/goshop/usecase/product"
	"github.com/gin-gonic/gin"
)

type productService interface {
	Create(ctx context.Context, input ucproduct.CreateProductInput) (*entity.Product, error)
	GetByID(ctx context.Context, id string) (*entity.Product, error)
	List(ctx context.Context, pagination valueobject.Pagination) ([]*entity.Product, error)
	Update(ctx context.Context, input ucproduct.UpdateProductInput) (*entity.Product, error)
	Delete(ctx context.Context, id string) error
	UploadImage(ctx context.Context, image ucproduct.UploadProductImage) (string, error)
}

type ProductHandler struct {
	service productService
}

func NewProductHandler(svc productService) *ProductHandler {
	return &ProductHandler{service: svc}
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req request.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	product, err := h.service.Create(c.Request.Context(), mapper.ToCreateProductInput(req))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, mapper.ToProductResponse(product))
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	product, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, mapper.ToProductResponse(product))
}

func (h *ProductHandler) List(c *gin.Context) {
	var req request.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}

	products, err := h.service.List(c.Request.Context(), req.ToPagination())
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, mapper.ToProductListResponse(products))
}

func (h *ProductHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req request.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	product, err := h.service.Update(c.Request.Context(), ucproduct.UpdateProductInput{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	})

	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, mapper.ToProductResponse(product))
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := h.service.Delete(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ProductHandler) UploadImage(c *gin.Context) {
	id := c.Param("id")
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.Error(err)
		return
	}

	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		c.Error(err)
		return
	}

	contentType := header.Header.Get("Content-Type")

	imageURL, err := h.service.UploadImage(c.Request.Context(), ucproduct.UploadProductImage{
		ProductID:   id,
		Data:        data,
		ContentType: contentType,
	})
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"image": imageURL})
}
