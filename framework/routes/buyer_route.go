package routes

import (
	"github.com/gofiber/fiber/v2"
	buyercontroller "github.com/hieronimusbudi/komodo-backend/controllers/buyer_controller"
	"github.com/hieronimusbudi/komodo-backend/dependencies"
	buyerrepo "github.com/hieronimusbudi/komodo-backend/framework/persistance/mysql/buyer_repository"
	buyerusecase "github.com/hieronimusbudi/komodo-backend/usecases/buyer_usecase"
)

// buyerRoutes used to define route and inject dependencies to repository, usecase and controller
func buyerRoutes(app *fiber.App, d *dependencies.Dependencies) {
	// inject connection to repository
	r := buyerrepo.NewMysqlBuyerRepository(d.Conn)
	// inject repository to usecase
	u := buyerusecase.NewBuyerUsecase(r)
	// inject usecase to controller
	c := buyercontroller.NewBuyerController(u)

	app.Post("/buyers/register", c.Register)
	app.Post("/buyers/login", c.Login)
}
