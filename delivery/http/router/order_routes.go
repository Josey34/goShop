package router

import (
	"github.com/Josey34/goshop/delivery/http/handler"
	"github.com/gin-gonic/gin"
)

func SetupOrderRoutes(rg *gin.RouterGroup, h *handler.OrderHandler) {
	rg.POST("", h.Create)
	rg.GET("", h.List)
	rg.GET("/:id", h.GetByID)
}
