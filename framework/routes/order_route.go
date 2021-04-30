package routes

import (
	"github.com/gofiber/fiber/v2"
	ordercontroller "github.com/hieronimusbudi/komodo-backend/controllers/order_controller"
	"github.com/hieronimusbudi/komodo-backend/dependencies"
	"github.com/hieronimusbudi/komodo-backend/framework/middlerwares"
	orderrepo "github.com/hieronimusbudi/komodo-backend/framework/persistance/mysql/order_repository"
	orderusecase "github.com/hieronimusbudi/komodo-backend/usecases/order_usecase"
)

// orderRoutes used to define route and inject dependencies to repository, usecase and controller
func orderRoutes(app *fiber.App, d *dependencies.Dependencies) {
	// inject connection to repository
	r := orderrepo.NewMysqlOrderRepository(d.Conn)
	// inject repository to usecase
	u := orderusecase.NewOrderUsecase(r)
	// inject usecase to controller
	c := ordercontroller.NewOrderController(u)

	app.Get("/orders/find/userid", middlerwares.ValidateRequest, c.GetByUserID)
	app.Post("/orders", middlerwares.ValidateRequest, middlerwares.BuyerTypeChecker, c.Store)
	app.Put("/orders/:id/accept", middlerwares.ValidateRequest, middlerwares.SellerTypeChecker, c.AcceptOrder)
}
