package routes

import (
	"github.com/gofiber/fiber/v2"
	ordercontroller "github.com/hieronimusbudi/komodo-backend/controllers/order_controller"
	productcontroller "github.com/hieronimusbudi/komodo-backend/controllers/product_controller"
	"github.com/hieronimusbudi/komodo-backend/dependencies"
	orderrepo "github.com/hieronimusbudi/komodo-backend/framework/persistance/mysql/order_repository"
	productrepo "github.com/hieronimusbudi/komodo-backend/framework/persistance/mysql/product_repository"
	orderusecase "github.com/hieronimusbudi/komodo-backend/usecases/order_usecase"
	productusecase "github.com/hieronimusbudi/komodo-backend/usecases/product_usecase"
)

// this function combines all routes and passes dependencies to routes
func All(app *fiber.App, d *dependencies.Dependencies) {

	// product
	rP := productrepo.NewMysqlProductRepository(d.Conn)
	uP := productusecase.NewProductUsecase(rP)
	cP := productcontroller.NewProductController(uP, d.Validate)

	// order
	rO := orderrepo.NewMysqlOrderRepository(d.Conn)
	uO := orderusecase.NewOrderUsecase(rO, rP)
	cO := ordercontroller.NewOrderController(uO, d.Validate)

	buyerRoutes(app, d)
	sellerRoutes(app, d)
	productRoutes(app, &cP)
	orderRoutes(app, &cO)
}
