package buyerrepo

import (
	"database/sql"
	"log"

	"github.com/hieronimusbudi/komodo-backend/entity"
)

const (
	queryGetAll  = "SELECT id, email, name, sending_address FROM buyers;"
	queryInsert  = "INSERT INTO buyers(email, name, password, sending_address) VALUES(?, ?, ?, ?);"
	queryGetById = "SELECT id, email, name, sending_address FROM buyers WHERE id=?;"
	queryUpdate  = "UPDATE buyers SET email=?, name=?, sending_address=? WHERE id=?;"
	queryDelete  = "DELETE FROM buyers WHERE id=?;"

	queryFindByEmail = "SELECT id, email, name, password, sending_address FROM buyers WHERE email=?;"
)

type mysqlBuyerRepository struct {
	Conn *sql.DB
}

func NewMysqlBuyerRepository(Conn *sql.DB) entity.BuyerRepository {
	return &mysqlBuyerRepository{Conn: Conn}
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
	if getErr := dbRes.Scan(&buyer.ID, &buyer.Email, &buyer.Name, &buyer.SendingAddress); getErr != nil {
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

	// email, name, password, sending_address
	dbRes, err := stmt.Exec(buyer.Email, buyer.Name, buyer.Password, buyer.SendingAddress)
	if err != nil {
		return err
	}

	buyerID, err := dbRes.LastInsertId()
	if err != nil {
		return err
	}
	log.Printf("%d \n", buyerID)
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

func (m *mysqlBuyerRepository) GetByEmail(buyer *entity.Buyer) (entity.Buyer, error) {
	stmt, err := m.Conn.Prepare(queryFindByEmail)
	if err != nil {
		return *buyer, err
	}
	defer stmt.Close()

	dbRes := stmt.QueryRow(buyer.Email)
	if getErr := dbRes.Scan(&buyer.ID, &buyer.Email, &buyer.Name, &buyer.Password, &buyer.SendingAddress); getErr != nil {
		return *buyer, getErr
	}
	return *buyer, nil
}
