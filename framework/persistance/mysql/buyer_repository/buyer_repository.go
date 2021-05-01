package buyerrepo

import (
	"database/sql"

	"github.com/hieronimusbudi/komodo-backend/entity"
	resterrors "github.com/hieronimusbudi/komodo-backend/framework/helpers/rest_errors"
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

// NewMysqlBuyerRepository will create a object with entity.BuyerRepository interface representation
func NewMysqlBuyerRepository(Conn *sql.DB) entity.BuyerRepository {
	return &mysqlBuyerRepository{Conn: Conn}
}

func (m *mysqlBuyerRepository) GetAll() ([]entity.Buyer, resterrors.RestErr) {
	dbRes, err := m.Conn.Query(queryGetAll)
	if err != nil {
		return nil, resterrors.NewInternalServerError("error when trying to get data", err)
	}
	defer m.Conn.Close()

	buyer := entity.Buyer{}
	res := []entity.Buyer{}
	for dbRes.Next() {
		var id int64
		var email, name, sendingAddress string
		err = dbRes.Scan(&id, &email, &name, &sendingAddress)
		if err != nil {
			return nil, resterrors.NewInternalServerError("error when trying to get data", err)
		}

		buyer.ID = id
		buyer.Name = name
		buyer.SendingAddress = sendingAddress

		res = append(res, buyer)
	}
	return res, nil
}

func (m *mysqlBuyerRepository) GetByID(buyer *entity.Buyer) resterrors.RestErr {
	stmt, err := m.Conn.Prepare(queryGetById)
	if err != nil {
		return resterrors.NewInternalServerError("error when trying to get data", err)
	}
	defer stmt.Close()

	dbRes := stmt.QueryRow(buyer.ID)
	if err := dbRes.Scan(&buyer.ID, &buyer.Email, &buyer.Name, &buyer.SendingAddress); err != nil {

		return resterrors.NewInternalServerError("error when trying to get data", err)
	}

	return nil
}

func (m *mysqlBuyerRepository) Store(buyer *entity.Buyer) resterrors.RestErr {
	stmt, err := m.Conn.Prepare(queryInsert)
	if err != nil {
		return resterrors.NewInternalServerError("error when trying to save data", err)
	}
	defer stmt.Close()

	// email, name, password, sending_address
	dbRes, err := stmt.Exec(buyer.Email, buyer.Name, buyer.Password, buyer.SendingAddress)
	if err != nil {
		return resterrors.NewInternalServerError("error when trying to save data", err)
	}

	buyerID, err := dbRes.LastInsertId()
	if err != nil {
		return resterrors.NewInternalServerError("error when trying to save data", err)
	}

	buyer.ID = buyerID

	return nil
}

func (m *mysqlBuyerRepository) Update(buyer *entity.Buyer) resterrors.RestErr {
	stmt, err := m.Conn.Prepare(queryUpdate)
	if err != nil {
		return resterrors.NewInternalServerError("error when trying to update data", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(buyer.Email, buyer.Name, buyer.SendingAddress, buyer.ID)
	if err != nil {
		return resterrors.NewInternalServerError("error when trying to update data", err)
	}
	return nil
}

func (m *mysqlBuyerRepository) Delete(buyer *entity.Buyer) resterrors.RestErr {
	stmt, err := m.Conn.Prepare(queryDelete)
	if err != nil {
		return resterrors.NewInternalServerError("error when trying to delete data", err)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(buyer.ID); err != nil {
		return resterrors.NewInternalServerError("error when trying to delete data", err)
	}
	return nil
}

func (m *mysqlBuyerRepository) GetByEmail(buyer *entity.Buyer) (entity.Buyer, resterrors.RestErr) {
	stmt, err := m.Conn.Prepare(queryFindByEmail)
	if err != nil {
		return *buyer, resterrors.NewInternalServerError("error when trying to get data", err)
	}
	defer stmt.Close()

	dbRes := stmt.QueryRow(buyer.Email)
	if err := dbRes.Scan(&buyer.ID, &buyer.Email, &buyer.Name, &buyer.Password, &buyer.SendingAddress); err != nil {
		return *buyer, resterrors.NewInternalServerError("error when trying to get data", err)
	}
	return *buyer, nil
}
