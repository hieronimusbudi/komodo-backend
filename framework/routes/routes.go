package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hieronimusbudi/komodo-backend/dependencies"
)

// this function combines all routes and passes dependencies to routes
func All(app *fiber.App, d *dependencies.Dependencies) {
	buyerRoutes(app, d)
	sellerRoutes(app, d)
	productRoutes(app, d)
	orderRoutes(app, d)
}
