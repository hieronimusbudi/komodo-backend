package orderrepo_test

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/hieronimusbudi/komodo-backend/entity"
	"github.com/hieronimusbudi/komodo-backend/framework/helpers"
	orderrepo "github.com/hieronimusbudi/komodo-backend/framework/persistance/mysql/order_repository"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type TestSuite struct {
	suite.Suite
	db                   *sql.DB
	mock                 sqlmock.Sqlmock
	repo                 entity.OrderRepository
	password             string
	hashedPassword       []byte
	expectedOrder1       entity.Order
	expectedOrderDetail1 entity.OrderDetail
	expectedProduct1     entity.Product
	expectedSeller1      entity.Seller
	expectedBuyer1       entity.Buyer
	price                []uint8
	time                 []uint8
}

// before each test
func (suite *TestSuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New()
	suite.NoError(err)

	// initiate repo
	suite.repo = orderrepo.NewMysqlOrderRepository(suite.db)

	// encrypt password
	suite.password = "12345"
	suite.hashedPassword, err = bcrypt.GenerateFromPassword([]byte(suite.password), bcrypt.DefaultCost)
	suite.NoError(err)

	suite.expectedBuyer1 = entity.Buyer{
		ID:             1,
		Email:          "buyer1@mail.com",
		Name:           "buyer",
		Password:       string(suite.hashedPassword),
		SendingAddress: "sending address",
	}

	suite.expectedSeller1 = entity.Seller{
		ID:            1,
		Email:         "seller1@mail.com",
		Name:          "seller",
		Password:      string(suite.hashedPassword),
		PickUpAddress: "pickup address",
	}

	suite.expectedProduct1 = entity.Product{
		ID:          1,
		Name:        "product1",
		Description: "desc",
		Price:       decimal.NewFromFloat(181818.11),
		Seller:      suite.expectedSeller1,
	}

	t, err := helpers.GetTimeNow()
	suite.NoError(err)
	suite.expectedOrder1 = entity.Order{
		ID:                         1,
		Buyer:                      suite.expectedBuyer1,
		Seller:                     suite.expectedSeller1,
		DeliverySourceAddress:      "pickup address",
		DeliveryDestinationAddress: "sending address",
		TotalQuantity:              10,
		TotalPrice:                 decimal.NewFromFloat(181818.11),
		Status:                     entity.PENDING,
		OrderDate:                  t,
		Items: []entity.OrderDetail{
			suite.expectedOrderDetail1,
		},
	}

	suite.expectedOrderDetail1 = entity.OrderDetail{
		ID:       1,
		Product:  suite.expectedProduct1,
		Quantity: 10,
	}
	suite.price = []uint8("181818.11")
	suite.time = []uint8(helpers.GetStringTimeNow())
}

func TestOrderRepo(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) TestGetByBuyerID() {
	queryGetByBuyerID := `SELECT id, buyer_id, seller_id, delivery_source_address, delivery_destination_address, 
	total_quantity, total_price, status, order_date FROM orders WHERE buyer_id=?;`
	prep := suite.mock.ExpectPrepare(regexp.QuoteMeta(queryGetByBuyerID))

	row1 := sqlmock.NewRows([]string{"id", "buyer_id", "seller_id", "delivery_source_address",
		"delivery_destination_address", "total_quantity", "total_price", "status", "order_date"}).
		AddRow(suite.expectedOrder1.ID, suite.expectedOrder1.Buyer.ID, suite.expectedOrder1.Seller.ID,
			suite.expectedOrder1.DeliverySourceAddress, suite.expectedOrder1.DeliveryDestinationAddress, suite.expectedOrder1.TotalQuantity,
			suite.price, suite.expectedOrder1.Status, suite.time,
		)
	prep.ExpectQuery().WillReturnRows(row1)

	odGetByOrderId := `SELECT id, product_id, quantity FROM order_details WHERE order_id=?;`
	expect := suite.mock.ExpectQuery(regexp.QuoteMeta(odGetByOrderId))
	row2 := sqlmock.NewRows([]string{"id", "product_id", "quantity"}).
		AddRow(suite.expectedOrderDetail1.ID, suite.expectedOrderDetail1.Product.ID, suite.expectedOrderDetail1.Quantity)
	expect.WillReturnRows(row2)

	res, repoErr := suite.repo.GetByBuyerID(suite.expectedOrder1.Buyer.ID)
	suite.NoError(repoErr)
	suite.NotNil(res)
}

func (suite *TestSuite) TestStore() {
	queryInsert := `INSERT INTO orders(buyer_id, seller_id, delivery_source_address, delivery_destination_address, 
		total_quantity, total_price, status, order_date) VALUES(?, ?, ?, ?, ?, ?, ?, ?);`
	odInsert := `INSERT INTO order_details(order_id, product_id, quantity) VALUES(?, ?, ?);`

	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(regexp.QuoteMeta(queryInsert)).
		WithArgs(suite.expectedOrder1.Buyer.ID, suite.expectedOrder1.Seller.ID, suite.expectedOrder1.DeliverySourceAddress,
			suite.expectedOrder1.DeliveryDestinationAddress, suite.expectedOrder1.TotalQuantity,
			suite.price, suite.expectedOrder1.Status, suite.time).
		WillReturnResult(sqlmock.NewResult(suite.expectedOrder1.ID, 1))

	suite.mock.ExpectExec(regexp.QuoteMeta(odInsert)).
		WithArgs(suite.expectedOrder1.ID, suite.expectedOrderDetail1.Product.ID, 10).
		WillReturnResult(sqlmock.NewResult(suite.expectedOrderDetail1.ID, 1))
	suite.mock.ExpectCommit()

	order := new(entity.Order)
	order.Buyer = suite.expectedBuyer1
	order.Seller = suite.expectedSeller1
	order.DeliverySourceAddress = suite.expectedOrder1.DeliverySourceAddress
	order.DeliveryDestinationAddress = suite.expectedOrder1.DeliveryDestinationAddress
	order.TotalQuantity = suite.expectedOrder1.TotalQuantity
	order.TotalPrice = suite.expectedOrder1.TotalPrice
	order.Status = suite.expectedOrder1.Status
	order.OrderDate = suite.expectedOrder1.OrderDate
	order.Items = []entity.OrderDetail{suite.expectedOrderDetail1}

	repoErr := suite.repo.Store(order)

	suite.NoError(repoErr)
	suite.NotNil(order)
}

func (suite *TestSuite) TestUpdate() {
	queryUpdate := `UPDATE orders SET buyer_id=?, seller_id=?, delivery_source_address=?, delivery_destination_address=?, 
	total_quantity=?, total_price=?, status=?, order_date=? WHERE id=?;`
	prep := suite.mock.ExpectPrepare(regexp.QuoteMeta(queryUpdate))

	prep.ExpectExec().
		WithArgs(suite.expectedOrder1.Buyer.ID, suite.expectedOrder1.Seller.ID, suite.expectedOrder1.DeliverySourceAddress,
			suite.expectedOrder1.DeliveryDestinationAddress, suite.expectedOrder1.TotalQuantity,
			suite.expectedOrder1.TotalPrice, suite.expectedOrder1.Status, suite.expectedOrder1.OrderDate, suite.expectedOrder1.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	order := new(entity.Order)
	order.ID = suite.expectedOrder1.ID
	order.Buyer = suite.expectedBuyer1
	order.Seller = suite.expectedSeller1
	order.DeliverySourceAddress = suite.expectedOrder1.DeliverySourceAddress
	order.DeliveryDestinationAddress = suite.expectedOrder1.DeliveryDestinationAddress
	order.TotalQuantity = suite.expectedOrder1.TotalQuantity
	order.TotalPrice = suite.expectedOrder1.TotalPrice
	order.Status = suite.expectedOrder1.Status
	order.OrderDate = suite.expectedOrder1.OrderDate
	order.Items = []entity.OrderDetail{suite.expectedOrderDetail1}

	repoErr := suite.repo.Update(order)
	suite.NoError(repoErr)
	suite.NotNil(order)
}
