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

func (d *Dependencies) Init() {
	d.Conn = mysqlpersistence.Client
	d.Validate = validator.New()
}
