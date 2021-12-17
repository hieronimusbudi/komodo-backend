// Integration Service
//
// This documentation describes example APIs found under https://github.com/ribice/golang-swaggerui-example
//
//     Schemes: https
//     BasePath: /v1
//     Version: 1.0.0
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - bearer
//
//     SecurityDefinitions:
//     bearer:
//          type: apiKey
//          name: Authorization
//          in: header
//
// swagger:meta
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
