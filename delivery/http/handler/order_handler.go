package handler

import (
	"net/http"

	"github.com/Josey34/goshop/delivery/http/dto/mapper"
	"github.com/Josey34/goshop/delivery/http/dto/request"
	"github.com/Josey34/goshop/delivery/http/dto/response"
	"github.com/Josey34/goshop/service"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	return &OrderHandler{service: svc}
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req request.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	customerID := c.GetString("customer_id")

	order, err := h.service.CreateOrder(c.Request.Context(), mapper.ToCreateOrderInput(customerID, req))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, mapper.ToOrderResponse(order))
}

func (h *OrderHandler) List(c *gin.Context) {
	var req request.PaginationRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}

	customerID := c.GetString("customer_id")
	orders, err := h.service.ListByCustomer(c.Request.Context(), customerID, req.ToPagination())

	if err != nil {
		c.Error(err)
		return
	}

	result := make([]response.OrderResponse, len(orders))
	for i, o := range orders {
		result[i] = mapper.ToOrderResponse(o)
	}
	c.JSON(http.StatusOK, result)
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	order, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, mapper.ToOrderResponse(order))
}
