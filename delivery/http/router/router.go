package router

import (
	"github.com/Josey34/goshop/delivery/http/handler"
	"github.com/Josey34/goshop/delivery/http/middleware"
	"github.com/Josey34/goshop/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	productH *handler.ProductHandler,
	orderH *handler.OrderHandler,
	authH *handler.AuthHandler,
	healthH *handler.HealthHandler,
	jwtSvc *jwt.JWTService,
) *gin.Engine {
	r := gin.New()

	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.ErrorMiddleware())

	r.GET("/health", healthH.Check)

	auth := r.Group("/auth")
	SetupAuthRoutes(auth, authH)

	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(jwtSvc))
	SetupProductRoutes(api.Group("/products"), productH)
	SetupOrderRoutes(api.Group("/orders"), orderH)

	return r
}
