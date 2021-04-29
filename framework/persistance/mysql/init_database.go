package mysqlpersistence

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hieronimusbudi/komodo-backend/config"
)

// Use this as DB Client connection
var (
	Client *sql.DB
)

func init() {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		config.MYSQL_USERNAME, config.MYSQL_PASSWORD, config.MYSQL_HOST, config.MYSQL_SCHEMA,
	)
	var err error
	Client, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Println(err)
	}

	if err = Client.Ping(); err != nil {
		log.Println(err)
	}

	log.Println("database succesfully configured")
}
