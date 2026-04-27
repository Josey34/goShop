package router

import (
	"github.com/Josey34/goshop/delivery/http/handler"
	"github.com/gin-gonic/gin"
)

func SetupProductRoutes(rg *gin.RouterGroup, h *handler.ProductHandler) {
	rg.POST("", h.Create)
	rg.GET("", h.List)
	rg.GET("/:id", h.GetByID)
	rg.PUT("/:id", h.Update)
	rg.DELETE("/:id", h.Delete)
}
