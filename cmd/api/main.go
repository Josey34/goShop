package main

import (
	"log"

	"github.com/Josey34/goshop/config"
	"github.com/Josey34/goshop/database"
	httpdelivery "github.com/Josey34/goshop/delivery/http"
	"github.com/Josey34/goshop/delivery/http/handler"
	"github.com/Josey34/goshop/delivery/http/router"
	"github.com/Josey34/goshop/pkg/hasher"
	"github.com/Josey34/goshop/pkg/idgen"
	"github.com/Josey34/goshop/pkg/jwt"
	"github.com/Josey34/goshop/repository/memory"
	"github.com/Josey34/goshop/service"
	"github.com/Josey34/goshop/usecase/auth"
	"github.com/Josey34/goshop/usecase/order"
	"github.com/Josey34/goshop/usecase/product"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewSQLiteDB(cfg.DB.Path)
	if err != nil {
		log.Fatal(err)
	}

	if err := database.RunMigrations(db, "database/migrations"); err != nil {
		log.Fatal(err)
	}

	idGen := idgen.NewUUIDGenerator()
	pwHasher := hasher.NewBcryptHasher()
	jwtSvc := jwt.NewJWTService(cfg.JWT.Secret, cfg.JWT.Expiry)

	productRepo := memory.NewProductRepo()
	orderRepo := memory.NewOrderRepo()
	customerRepo := memory.NewCustomerRepo()
	orderQueue := memory.NewOrderQueue()

	createProductUC := product.NewCreateProductUseCase(productRepo, idGen)
	getProductUC := product.NewGetProductUseCase(productRepo)
	listProductUC := product.NewListProductUseCase(productRepo)
	updateProductUC := product.NewUpdateProductUseCase(productRepo)
	deleteProductUC := product.NewDeleteProductUseCase(productRepo)

	createOrderUC := order.NewCreateOrderUseCase(orderRepo, productRepo, orderQueue, idGen)
	getOrderUC := order.NewGetOrderUseCase(orderRepo)
	listOrderUC := order.NewListProductUseCase(orderRepo)
	updateOrderUC := order.NewUpdateOrderStatusUseCase(orderRepo)

	loginUC := auth.NewLoginUseCase(customerRepo, pwHasher)
	registerUC := auth.NewRegisterUseCase(customerRepo, pwHasher, idGen)

	productSvc := service.NewProductService(createProductUC, getProductUC, listProductUC, updateProductUC, deleteProductUC)
	orderSvc := service.NewOrderService(createOrderUC, getOrderUC, listOrderUC, updateOrderUC)
	authSvc := service.NewAuthService(registerUC, loginUC, jwtSvc)

	productH := handler.NewProductHandler(productSvc)
	orderH := handler.NewOrderHandler(orderSvc)
	authH := handler.NewAuthHandler(authSvc)
	healthH := handler.NewHealthHandler()

	engine := router.SetupRouter(productH, orderH, authH, healthH, jwtSvc)

	srv := httpdelivery.NewServer(engine, cfg.App.Port)
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
