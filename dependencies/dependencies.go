package dependencies

import (
	"database/sql"

	"github.com/go-playground/validator/v10"
	mysqlpersistence "github.com/hieronimusbudi/komodo-backend/framework/persistance/mysql"
)

type Dependencies struct {
	Conn     *sql.DB
	Validate *validator.Validate
}

func NewDependencies() *Dependencies {
	conn := mysqlpersistence.Client
	validate := validator.New()
	return &Dependencies{
		Conn:     conn,
		Validate: validate,
	}
}
