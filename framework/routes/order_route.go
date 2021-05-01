package routes

import (
	"github.com/gofiber/fiber/v2"
	ordercontroller "github.com/hieronimusbudi/komodo-backend/controllers/order_controller"
	"github.com/hieronimusbudi/komodo-backend/framework/middlerwares"
)

// orderRoutes used to define route and inject dependencies to repository, usecase and controller
func orderRoutes(app *fiber.App, c *ordercontroller.OrderController) {
	app.Get("/orders/find/byuser", middlerwares.ValidateRequest, (*c).GetByUserID)
	app.Post("/orders", middlerwares.ValidateRequest, middlerwares.BuyerTypeChecker, (*c).Store)
	app.Put("/orders/:id/accept", middlerwares.ValidateRequest, middlerwares.SellerTypeChecker, (*c).AcceptOrder)
}
