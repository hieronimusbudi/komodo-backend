package orderrepo

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/hieronimusbudi/komodo-backend/entity"
	"github.com/shopspring/decimal"
)

const (
	queryGetAll = `SELECT id, buyer_id, seller_id, delivery_source_address, delivery_destination_address, 
	total_quantity, total_price, status, order_date 
	FROM orders;`
	queryGetById = `SELECT id, buyer_id, seller_id, delivery_source_address, delivery_destination_address, 
	total_quantity, total_price, status, order_date FROM orders WHERE id=?;`
	queryGetByBuyerID = `SELECT id, buyer_id, seller_id, delivery_source_address, delivery_destination_address, 
	total_quantity, total_price, status, order_date FROM orders WHERE buyer_id=?;`
	queryGetBySellerID = `SELECT id, buyer_id, seller_id, delivery_source_address, delivery_destination_address, 
	total_quantity, total_price, status, order_date FROM orders WHERE seller_id=?;`

	queryInsert = `INSERT INTO orders(buyer_id, seller_id, delivery_source_address, delivery_destination_address, 
		total_quantity, total_price, status, order_date) VALUES(?, ?, ?, ?, ?, ?, ?, ?);`
	queryUpdate = `UPDATE orders SET buyer_id=?, seller_id=?, delivery_source_address=?, delivery_destination_address=?, 
	total_quantity=?, total_price=?, status=?, order_date=? WHERE id=?;`
	queryDelete = "DELETE FROM orders WHERE id=?;"

	odInsert       = `INSERT INTO order_details(order_id, product_id, quantity) VALUES(?, ?, ?);`
	odGetByOrderId = `SELECT id, product_id, quantity FROM order_details WHERE order_id=?;`
)

type mysqlOrderRepository struct {
	Conn *sql.DB
}

func NewMysqlOrderRepository(Conn *sql.DB) entity.OrderRepository {
	return &mysqlOrderRepository{Conn: Conn}
}

func (m *mysqlOrderRepository) GetAll() ([]entity.Order, error) {
	stmt, err := m.Conn.Prepare(queryGetAll)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	dbRes, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	order := entity.Order{}
	res := []entity.Order{}
	for dbRes.Next() {
		var totalPrice []uint8

		err = dbRes.Scan(
			&order.ID, &order.Buyer.ID, &order.Seller.ID, &order.DeliverySourceAddress,
			&order.DeliveryDestinationAddress, &order.TotalQuantity, totalPrice, &order.Status, &order.OrderDate)
		if err != nil {
			panic(err.Error())
		}

		dP, err := decimal.NewFromString(string(totalPrice))
		if err != nil {
			panic(err.Error())
		}
		order.TotalPrice = dP

		res = append(res, order)
	}
	return res, nil
}

func (m *mysqlOrderRepository) GetByBuyerID(buyerID int64) ([]entity.Order, error) {
	stmt, err := m.Conn.Prepare(queryGetByBuyerID)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	dbRes, err := stmt.Query(buyerID)
	if err != nil {
		return nil, err
	}

	orderRow := entity.Order{}
	res := []entity.Order{}
	for dbRes.Next() {
		var totalPrice, orderDate []uint8

		// get oder for each row
		err = dbRes.Scan(
			&orderRow.ID, &orderRow.Buyer.ID, &orderRow.Seller.ID, &orderRow.DeliverySourceAddress,
			&orderRow.DeliveryDestinationAddress, &orderRow.TotalQuantity, &totalPrice, &orderRow.Status, &orderDate)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}

		dP, err := decimal.NewFromString(string(totalPrice))
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		orderRow.TotalPrice = dP

		vT, err := time.Parse("2006-01-02 15:04:05", string(orderDate))
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		orderRow.OrderDate = vT

		// find order detail for each order
		odRes, err := m.Conn.Query(odGetByOrderId, orderRow.ID)
		if err != nil {
			return nil, err
		}

		// scan order details result
		odRow := entity.OrderDetail{}
		orderRow.Items = []entity.OrderDetail{}
		for odRes.Next() {
			// id, order_id, product_id, quantity
			err = odRes.Scan(&odRow.ID, &odRow.Product.ID, &odRow.Quantity)
			if err != nil {
				log.Println(112, err.Error())
				return nil, err
			}

			orderRow.Items = append(orderRow.Items, odRow)
		}

		res = append(res, orderRow)
	}
	return res, nil
}

func (m *mysqlOrderRepository) GetBySellerID(sellerID int64) ([]entity.Order, error) {
	stmt, err := m.Conn.Prepare(queryGetBySellerID)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	dbRes, err := stmt.Query(sellerID)
	if err != nil {
		return nil, err
	}

	orderRow := entity.Order{}
	res := []entity.Order{}
	for dbRes.Next() {
		var totalPrice, orderDate []uint8

		// get oder for each row
		err = dbRes.Scan(
			&orderRow.ID, &orderRow.Buyer.ID, &orderRow.Seller.ID, &orderRow.DeliverySourceAddress,
			&orderRow.DeliveryDestinationAddress, &orderRow.TotalQuantity, &totalPrice, &orderRow.Status, &orderDate)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}

		dP, err := decimal.NewFromString(string(totalPrice))
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		orderRow.TotalPrice = dP

		vT, err := time.Parse("2006-01-02 15:04:05", string(orderDate))
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		orderRow.OrderDate = vT

		// find order detail for each order
		odRes, err := m.Conn.Query(odGetByOrderId, orderRow.ID)
		if err != nil {
			return nil, err
		}

		// scan order details result
		odRow := entity.OrderDetail{}
		orderRow.Items = []entity.OrderDetail{}
		for odRes.Next() {
			// id, order_id, product_id, quantity
			err = odRes.Scan(&odRow.ID, &odRow.Product.ID, &odRow.Quantity)
			if err != nil {
				log.Println(112, err.Error())
				return nil, err
			}

			orderRow.Items = append(orderRow.Items, odRow)
		}

		res = append(res, orderRow)
	}
	return res, nil
}

func (m *mysqlOrderRepository) GetByID(order *entity.Order) error {
	stmt, err := m.Conn.Prepare(queryGetById)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var totalPrice, orderDate []uint8
	dbRes := stmt.QueryRow(order.ID)
	if getErr := dbRes.Scan(&order.ID, &order.Buyer.ID, &order.Seller.ID, &order.DeliverySourceAddress,
		&order.DeliveryDestinationAddress, &order.TotalQuantity, &totalPrice, &order.Status, &orderDate); getErr != nil {
		return getErr
	}

	dP, err := decimal.NewFromString(string(totalPrice))
	if err != nil {
		return err
	}
	order.TotalPrice = dP

	vT, err := time.Parse("2006-01-02 15:04:05", string(orderDate))
	if err != nil {
		return err
	}
	order.OrderDate = vT

	return nil
}

func (m *mysqlOrderRepository) Store(order *entity.Order) error {
	// start transaction sequence
	ctx := context.Background()
	tx, err := m.Conn.BeginTx(ctx, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	// insert order
	dbRes, err := tx.ExecContext(
		ctx, queryInsert,
		order.Buyer.ID, order.Seller.ID, order.DeliverySourceAddress, order.DeliveryDestinationAddress,
		order.TotalQuantity, []uint8(order.TotalPrice.String()), order.Status, []uint8(order.OrderDate.Format("2006-01-02 15:04:05")))
	if err != nil {
		log.Println(1, err)
		tx.Rollback()
		return err
	}

	orderID, err := dbRes.LastInsertId()
	if err != nil {
		log.Println(2, err)
		tx.Rollback()
		return err
	}
	order.ID = orderID

	// insert order details
	for idx, od := range order.Items {
		odRes, err := tx.ExecContext(
			ctx, odInsert,
			orderID, od.Product.ID, od.Quantity)
		if err != nil {
			tx.Rollback()
			return err
		}

		odID, err := odRes.LastInsertId()
		if err != nil {
			tx.Rollback()
			return err
		}
		order.Items[idx].ID = odID
	}

	// commit the change if all queries ran successfully
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (m *mysqlOrderRepository) Update(order *entity.Order) error {
	stmt, err := m.Conn.Prepare(queryUpdate)
	if err != nil {
		return err
	}
	defer stmt.Close()

	//buyer_id, seller_id, delivery_source_address, delivery_destination_address,total_quantity, total_price, status, order_date
	_, err = stmt.Exec(order.Buyer.ID, order.Seller.ID, order.DeliverySourceAddress,
		order.DeliveryDestinationAddress, order.TotalQuantity, order.TotalPrice, order.Status, order.OrderDate, order.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m *mysqlOrderRepository) Delete(order *entity.Order) error {
	stmt, err := m.Conn.Prepare(queryDelete)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(order.ID); err != nil {
		return err
	}
	return nil
}
