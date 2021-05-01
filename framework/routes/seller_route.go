package routes

import (
	"github.com/gofiber/fiber/v2"
	sellercontroller "github.com/hieronimusbudi/komodo-backend/controllers/seller_controller"
	"github.com/hieronimusbudi/komodo-backend/dependencies"
	sellerrepo "github.com/hieronimusbudi/komodo-backend/framework/persistance/mysql/seller_repository"
	sellerusecase "github.com/hieronimusbudi/komodo-backend/usecases/seller_usecase"
)

// sellerRoutes used to define route and inject dependencies to repository, usecase and controller
func sellerRoutes(app *fiber.App, d *dependencies.Dependencies) {
	// inject connection to repository
	r := sellerrepo.NewMysqlSellerRepository(d.Conn)
	// inject repository to usecase
	u := sellerusecase.NewSellerUsecase(r)
	// inject usecase to controller
	c := sellercontroller.NewSellerController(u, d.Validate)

	app.Post("/sellers/register", c.Register)
	app.Post("/sellers/login", c.Login)
}
