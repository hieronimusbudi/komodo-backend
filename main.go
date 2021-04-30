package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hieronimusbudi/komodo-backend/dependencies"
	"github.com/hieronimusbudi/komodo-backend/framework/routes"
)

func main() {
	app := fiber.New()
	dependencies := new(dependencies.Dependencies)
	dependencies.Init()

	routes.All(app, dependencies)
	app.Listen(":9000")
}
