package repository

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hieronimusbudi/komodo-backend/entity"
)

const (
	queryGetAll  = "SELECT id, email, name, sending_address FROM buyers;"
	queryInsert  = "INSERT INTO buyers(email, name, password, sending_address) VALUES(?, ?, ?, ?, ?);"
	queryGetById = "SELECT id, email, name, sending_address FROM buyers WHERE id=?;"
	queryUpdate  = "UPDATE buyers SET email=?, name=?, sending_address=? WHERE id=?;"
	queryDelete  = "DELETE FROM buyers WHERE id=?;"

	queryFindByEmailAndPassword = "SELECT id, email, name, sending_address FROM buyers WHERE email=? AND password=?;"
)

type mysqlBuyerRepository struct {
	Conn *sql.DB
}

func NewMysqlBuyerRepository() entity.BuyerRepository {
	return &mysqlBuyerRepository{}
}

func (m *mysqlBuyerRepository) Init() {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		username, password, host, schema,
	)
	var err error
	m.Conn, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Println(err)
	}

	if err = m.Conn.Ping(); err != nil {
		log.Println(err)
	}

	log.Println("database succesfully configured")
}

func (m *mysqlBuyerRepository) GetAll() ([]entity.Buyer, error) {
	dbRes, err := m.Conn.Query(queryGetAll)
	if err != nil {
		return nil, err
	}
	defer m.Conn.Close()

	buyer := entity.Buyer{}
	res := []entity.Buyer{}
	for dbRes.Next() {
		var id int64
		var email, name, sendingAddress string
		err = dbRes.Scan(&id, &email, &name, &sendingAddress)
		if err != nil {
			panic(err.Error())
		}

		buyer.ID = id
		buyer.Name = name
		buyer.SendingAddress = sendingAddress

		res = append(res, buyer)
	}
	return res, nil
}

func (m *mysqlBuyerRepository) GetByID(buyer *entity.Buyer) error {
	stmt, err := m.Conn.Prepare(queryGetById)
	if err != nil {
		return err
	}
	defer stmt.Close()

	dbRes := stmt.QueryRow(buyer.ID)

	if getErr := dbRes.Scan(buyer.ID, buyer.Email, buyer.Name, buyer.SendingAddress); getErr != nil {
		return getErr
	}
	return nil
}

func (m *mysqlBuyerRepository) Store(buyer *entity.Buyer) error {
	stmt, err := m.Conn.Prepare(queryInsert)
	if err != nil {
		return err
	}
	defer stmt.Close()

	dbRes, err := stmt.Exec(buyer.Email, buyer.Name, buyer.Password, buyer.SendingAddress)
	if err != nil {
		return err
	}

	buyerID, err := dbRes.LastInsertId()
	if err != nil {
		return err
	}
	buyer.ID = buyerID

	return nil
}

func (m *mysqlBuyerRepository) Update(buyer *entity.Buyer) error {
	stmt, err := m.Conn.Prepare(queryUpdate)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(buyer.Email, buyer.Name, buyer.SendingAddress, buyer.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m *mysqlBuyerRepository) Delete(buyer *entity.Buyer) error {
	stmt, err := m.Conn.Prepare(queryDelete)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(buyer.ID); err != nil {
		return err
	}
	return nil
}

func (m *mysqlBuyerRepository) GetByEmailAndPassword(buyer *entity.Buyer) error {
	stmt, err := m.Conn.Prepare(queryFindByEmailAndPassword)
	if err != nil {
		return err
	}
	defer stmt.Close()

	dbRes := stmt.QueryRow(buyer.Email, buyer.Password)

	if getErr := dbRes.Scan(buyer.ID, buyer.Email, buyer.Name, buyer.SendingAddress); getErr != nil {
		return getErr
	}
	return nil
}
