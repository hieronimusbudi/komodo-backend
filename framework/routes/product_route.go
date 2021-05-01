package routes

import (
	"github.com/gofiber/fiber/v2"
	productcontroller "github.com/hieronimusbudi/komodo-backend/controllers/product_controller"
	"github.com/hieronimusbudi/komodo-backend/framework/middlerwares"
)

// productRoutes used to define route and inject dependencies to repository, usecase and controller
func productRoutes(app *fiber.App, c *productcontroller.ProductController) {
	app.Get("/products", (*c).GetAll)
	app.Post("/products", middlerwares.ValidateRequest, middlerwares.SellerTypeChecker, (*c).Store)
}
