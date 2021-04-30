package productrepo

import (
	"database/sql"

	"github.com/hieronimusbudi/komodo-backend/entity"
	resterrors "github.com/hieronimusbudi/komodo-backend/framework/helpers/rest_errors"
	"github.com/shopspring/decimal"
)

const (
	queryGetAll  = "SELECT id, name, description, price, seller_id FROM products;"
	queryInsert  = "INSERT INTO products(name, description, price, seller_id) VALUES(?, ?, ?, ?);"
	queryGetById = "SELECT id, name, description, price, seller_id FROM products WHERE id=?;"
	queryUpdate  = "UPDATE products SET name=?, description=?, price=?, seller_id=? WHERE id=?;"
	queryDelete  = "DELETE FROM products WHERE id=?;"
)

type mysqlProductRepository struct {
	Conn *sql.DB
}

// NewMysqlProductRepository will create a object with entity.ProductRepository interface representation
func NewMysqlProductRepository(Conn *sql.DB) entity.ProductRepository {
	return &mysqlProductRepository{Conn: Conn}
}

func (m *mysqlProductRepository) GetAll() ([]entity.Product, resterrors.RestErr) {
	stmt, err := m.Conn.Prepare(queryGetAll)
	if err != nil {
		return nil, resterrors.NewInternalServerError("error when trying to get data", err)
	}
	defer stmt.Close()

	dbRes, err := stmt.Query()
	if err != nil {
		return nil, resterrors.NewInternalServerError("error when trying to get data", err)
	}

	product := entity.Product{}
	res := []entity.Product{}
	for dbRes.Next() {
		var id, seller_id int64
		var price []uint8
		var name, description string
		err = dbRes.Scan(&id, &name, &description, &price, &seller_id)
		if err != nil {
			return nil, resterrors.NewInternalServerError("error when trying to get data", err)
		}

		product.ID = id
		product.Name = name
		product.Description = description
		product.Seller.ID = seller_id

		dP, err := decimal.NewFromString(string(price))
		if err != nil {
			return nil, resterrors.NewInternalServerError("error when trying to get data", err)
		}
		product.Price = dP

		res = append(res, product)
	}
	return res, nil
}

func (m *mysqlProductRepository) GetByID(product *entity.Product) resterrors.RestErr {
	stmt, err := m.Conn.Prepare(queryGetById)
	if err != nil {
		return resterrors.NewInternalServerError("error when trying to get data", err)
	}
	defer stmt.Close()

	var price []uint8
	dbRes := stmt.QueryRow(product.ID)
	if err := dbRes.Scan(&product.ID, &product.Name, &product.Description, &price, &product.Seller.ID); err != nil {
		return resterrors.NewInternalServerError("error when trying to get data", err)
	}

	dP, err := decimal.NewFromString(string(price))
	if err != nil {
		return resterrors.NewInternalServerError("error when trying to get data", err)
	}
	product.Price = dP

	return nil
}

func (m *mysqlProductRepository) Store(product *entity.Product) resterrors.RestErr {
	stmt, err := m.Conn.Prepare(queryInsert)
	if err != nil {
		return resterrors.NewInternalServerError("error when trying to save data", err)
	}
	defer stmt.Close()

	// name, description, price, seller_id
	dbRes, err := stmt.Exec(product.Name, product.Description, product.Price, product.Seller.ID)
	if err != nil {
		return resterrors.NewInternalServerError("error when trying to save data", err)
	}

	productID, err := dbRes.LastInsertId()
	if err != nil {
		return resterrors.NewInternalServerError("error when trying to save data", err)
	}

	product.ID = productID
	return nil
}

func (m *mysqlProductRepository) Update(product *entity.Product) resterrors.RestErr {
	stmt, err := m.Conn.Prepare(queryUpdate)
	if err != nil {
		return resterrors.NewInternalServerError("error when trying to update data", err)
	}
	defer stmt.Close()

	// name, description, price, seller_id
	_, err = stmt.Exec(product.Name, product.Description, product.Price, product.Seller.ID)
	if err != nil {
		return resterrors.NewInternalServerError("error when trying to update data", err)
	}
	return nil
}

func (m *mysqlProductRepository) Delete(product *entity.Product) resterrors.RestErr {
	stmt, err := m.Conn.Prepare(queryDelete)
	if err != nil {
		return resterrors.NewInternalServerError("error when trying to delete data", err)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(product.ID); err != nil {
		return resterrors.NewInternalServerError("error when trying to delete data", err)
	}
	return nil
}
