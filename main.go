package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hieronimusbudi/komodo-backend/framework/routes"
)

func main() {
	app := fiber.New()
	routes.All(app)
	app.Listen(":9000")
}
