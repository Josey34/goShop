package testutil

import (
	"testing"

	"github.com/Josey34/goshop/delivery/http/handler"
	"github.com/Josey34/goshop/delivery/http/router"
	"github.com/Josey34/goshop/pkg/hasher"
	"github.com/Josey34/goshop/pkg/idgen"
	"github.com/Josey34/goshop/pkg/jwt"
	"github.com/Josey34/goshop/repository/memory"
	"github.com/Josey34/goshop/service"
	ucauth "github.com/Josey34/goshop/usecase/auth"
	ucorder "github.com/Josey34/goshop/usecase/order"
	ucproduct "github.com/Josey34/goshop/usecase/product"
	"github.com/gin-gonic/gin"
)

const TestJWTSecret = "integration-test-secret"

func NewTestEngine(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	idGen := idgen.NewUUIDGenerator()
	pwHasher := hasher.NewBcryptHasher()
	jwtSvc := jwt.NewJWTService(TestJWTSecret, 1)

	productRepo := memory.NewProductRepo()
	orderRepo := memory.NewOrderRepo()
	customerRepo := memory.NewCustomerRepo()
	queue := memory.NewOrderQueue()
	fileStorage := memory.NewFileStorage()

	createProductUC := ucproduct.NewCreateProductUseCase(productRepo, idGen)
	getProductUC := ucproduct.NewGetProductUseCase(productRepo)
	listProductUC := ucproduct.NewListProductUseCase(productRepo)
	updateProductUC := ucproduct.NewUpdateProductUseCase(productRepo)
	deleteProductUC := ucproduct.NewDeleteProductUseCase(productRepo)
	uploadImageUC := ucproduct.NewUploadProductImageUsecase(productRepo, fileStorage)

	createOrderUC := ucorder.NewCreateOrderUseCase(orderRepo, productRepo, idGen)
	getOrderUC := ucorder.NewGetOrderUseCase(orderRepo)
	listOrderUC := ucorder.NewListProductUseCase(orderRepo)
	updateOrderUC := ucorder.NewUpdateOrderStatusUseCase(orderRepo)

	registerUC := ucauth.NewRegisterUseCase(customerRepo, pwHasher, idGen)
	loginUC := ucauth.NewLoginUseCase(customerRepo, pwHasher)

	productSvc := service.NewProductService(createProductUC, getProductUC, listProductUC, updateProductUC, deleteProductUC, uploadImageUC)
	orderSvc := service.NewOrderService(createOrderUC, getOrderUC, listOrderUC, updateOrderUC, queue)
	authSvc := service.NewAuthService(registerUC, loginUC, jwtSvc)

	productH := handler.NewProductHandler(productSvc)
	orderH := handler.NewOrderHandler(orderSvc)
	authH := handler.NewAuthHandler(authSvc)
	healthH := handler.NewHealthHandler()

	return router.SetupRouter(productH, orderH, authH, healthH, jwtSvc)
}
