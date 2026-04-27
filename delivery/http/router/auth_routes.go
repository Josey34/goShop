package router

import (
	"github.com/Josey34/goshop/delivery/http/handler"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(rg *gin.RouterGroup, h *handler.AuthHandler) {
	rg.POST("/login", h.Login)
	rg.POST("/register", h.Register)
}
