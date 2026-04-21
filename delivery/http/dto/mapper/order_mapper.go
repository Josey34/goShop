package mapper

import (
	"github.com/Josey34/goshop/delivery/http/dto/request"
	"github.com/Josey34/goshop/delivery/http/dto/response"
	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/factory"
	ucorder "github.com/Josey34/goshop/usecase/order"
)

func ToCreateOrderInput(customerID string, req request.CreateOrderRequest) ucorder.CreateOrderInput {
	result := make([]factory.CreateOrderItemInput, len(req.Items))
	for i, item := range req.Items {
		result[i] = factory.CreateOrderItemInput{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	return ucorder.CreateOrderInput{
		CustomerID: customerID,
		Items:      result,
	}
}

func ToOrderItemResponse(item entity.OrderItem) response.OrderItemResponse {
	return response.OrderItemResponse{
		ID:        item.ID,
		ProductID: item.ProductID,
		Name:      item.Name,
		Price:     item.Price.Value(),
		Quantity:  item.Quantity,
		LineTotal: item.LineTotal().Value(),
	}
}

func ToOrderResponse(o *entity.Order) response.OrderResponse {
	items := make([]response.OrderItemResponse, len(o.Items))
	for i, item := range o.Items {
		items[i] = ToOrderItemResponse(item)
	}

	return response.OrderResponse{
		ID:         o.ID,
		CustomerID: o.CustomerID,
		Items:      items,
		Total:      o.Total.Value(),
		Status:     string(o.Status),
		CreatedAt:  o.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
