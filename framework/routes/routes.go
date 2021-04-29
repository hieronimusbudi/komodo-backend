package routes

import (
	"github.com/gofiber/fiber/v2"
	mysqlpersistence "github.com/hieronimusbudi/komodo-backend/framework/persistance/mysql"
)

func All(app *fiber.App) {
	conn := mysqlpersistence.Client
	buyerRoutes(app, conn)
	sellerRoutes(app, conn)
	productRoutes(app, conn)
	orderRoutes(app, conn)
}
