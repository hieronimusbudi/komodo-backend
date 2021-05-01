package config

import "os"

var (
	PORT           = os.Getenv("PORT")
	JWT_SECRET     = os.Getenv("JWT_SECRET")
	MYSQL_USER     = os.Getenv("MYSQL_USER")
	MYSQL_PASSWORD = os.Getenv("MYSQL_PASSWORD")
	MYSQL_HOST     = os.Getenv("MYSQL_HOST")
	MYSQL_PORT     = os.Getenv("MYSQL_PORT")
	MYSQL_DATABASE = os.Getenv("MYSQL_DATABASE")
)
