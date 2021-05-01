package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/hieronimusbudi/komodo-backend/config"
	"github.com/hieronimusbudi/komodo-backend/dependencies"
	"github.com/hieronimusbudi/komodo-backend/framework/routes"
)

func main() {
	app := fiber.New()
	dependencies := dependencies.NewDependencies()

	routes.All(app, dependencies)
	app.Listen(fmt.Sprintf(":%s", config.PORT))
}
