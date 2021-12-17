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
	c := buyercontroller.NewBuyerController(u, d.Validate)

	app.Post("/buyers/register", c.Register)

	// swagger:operation POST /auth/login login loginRequest
	// ---
	// summary: Login.
	// description: Login using API key and Password.
	// parameters:
	// - in: body
	//   name: body
	//   description: Login body request
	//   required: true
	//   schema:
	//     $ref: "#/definitions/loginRequest"
	// responses:
	//   "200":
	//     "$ref": "#/definitions/loginRequest"
	app.Post("/buyers/login", c.Login)
}
