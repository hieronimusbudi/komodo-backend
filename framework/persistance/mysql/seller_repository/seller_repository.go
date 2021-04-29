package sellerrepo

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hieronimusbudi/komodo-backend/entity"
)

const (
	queryGetAll  = "SELECT id, email, name, pickup_address FROM sellers;"
	queryInsert  = "INSERT INTO sellers(email, name, password, pickup_address) VALUES(?, ?, ?, ?);"
	queryGetById = "SELECT id, email, name, pickup_address FROM sellers WHERE id=?;"
	queryUpdate  = "UPDATE sellers SET email=?, name=?, pickup_address=? WHERE id=?;"
	queryDelete  = "DELETE FROM sellers WHERE id=?;"

	queryFindByEmail = "SELECT id, email, name, password, pickup_address FROM sellers WHERE email=?;"
)

type mysqlSellerRepository struct {
	Conn *sql.DB
}

func NewMysqlSellerRepository(Conn *sql.DB) entity.SellerRepository {
	return &mysqlSellerRepository{Conn: Conn}
}

func (m *mysqlSellerRepository) GetAll() ([]entity.Seller, error) {
	dbRes, err := m.Conn.Query(queryGetAll)
	if err != nil {
		return nil, err
	}
	defer m.Conn.Close()

	seller := entity.Seller{}
	res := []entity.Seller{}
	for dbRes.Next() {
		var id int64
		var email, name, pickupAddress string
		err = dbRes.Scan(&id, &email, &name, &pickupAddress)
		if err != nil {
			panic(err.Error())
		}

		seller.ID = id
		seller.Name = name
		seller.PickUpAddress = pickupAddress

		res = append(res, seller)
	}
	return res, nil
}

func (m *mysqlSellerRepository) GetByID(seller *entity.Seller) error {
	stmt, err := m.Conn.Prepare(queryGetById)
	if err != nil {
		return err
	}
	defer stmt.Close()

	dbRes := stmt.QueryRow(seller.ID)

	if getErr := dbRes.Scan(&seller.ID, &seller.Email, &seller.Name, &seller.PickUpAddress); getErr != nil {
		return getErr
	}
	return nil
}

func (m *mysqlSellerRepository) Store(seller *entity.Seller) error {
	stmt, err := m.Conn.Prepare(queryInsert)
	if err != nil {
		return err
	}
	defer stmt.Close()
	// email, name, password, pickup_address
	dbRes, err := stmt.Exec(seller.Email, seller.Name, seller.Password, seller.PickUpAddress)
	if err != nil {
		return err
	}

	sellerID, err := dbRes.LastInsertId()
	if err != nil {
		return err
	}
	log.Printf("%d \n", sellerID)
	seller.ID = sellerID

	return nil
}

func (m *mysqlSellerRepository) Update(seller *entity.Seller) error {
	stmt, err := m.Conn.Prepare(queryUpdate)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(seller.Email, seller.Name, seller.PickUpAddress, seller.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m *mysqlSellerRepository) Delete(seller *entity.Seller) error {
	stmt, err := m.Conn.Prepare(queryDelete)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(seller.ID); err != nil {
		return err
	}
	return nil
}

func (m *mysqlSellerRepository) GetByEmail(seller *entity.Seller) (entity.Seller, error) {
	stmt, err := m.Conn.Prepare(queryFindByEmail)
	if err != nil {
		return *seller, err
	}
	defer stmt.Close()

	dbRes := stmt.QueryRow(seller.Email)

	if getErr := dbRes.Scan(&seller.ID, &seller.Email, &seller.Name, &seller.Password, &seller.PickUpAddress); getErr != nil {
		return *seller, getErr
	}
	return *seller, nil
}
