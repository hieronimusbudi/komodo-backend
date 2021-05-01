package productrepo_test

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/hieronimusbudi/komodo-backend/entity"
	productrepo "github.com/hieronimusbudi/komodo-backend/framework/persistance/mysql/product_repository"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type TestSuite struct {
	suite.Suite
	db               *sql.DB
	mock             sqlmock.Sqlmock
	repo             entity.ProductRepository
	password         string
	hashedPassword   []byte
	expectedProduct1 entity.Product
	expectedProduct2 entity.Product
	expectedSeller1  entity.Seller
	price            []uint8
}

// before each test
func (suite *TestSuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New()
	suite.NoError(err)
	// initiate repo
	suite.repo = productrepo.NewMysqlProductRepository(suite.db)

	// encrypt password
	suite.password = "12345"
	suite.hashedPassword, err = bcrypt.GenerateFromPassword([]byte(suite.password), bcrypt.DefaultCost)
	suite.NoError(err)

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

	suite.expectedProduct2 = entity.Product{
		ID:          2,
		Name:        "product1",
		Description: "desc",
		Price:       decimal.NewFromFloat(181818.11),
		Seller:      suite.expectedSeller1,
	}

	suite.price = []uint8("181818.11")
}

func TestProductRepo(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) TestGetAll() {
	queryGetAll := "SELECT id, name, description, price, seller_id FROM products;"
	prep := suite.mock.ExpectPrepare(regexp.QuoteMeta(queryGetAll))

	row1 := sqlmock.NewRows([]string{"id", "name", "description", "price", "seller_id"}).
		AddRow(suite.expectedProduct1.ID, suite.expectedProduct1.Name, suite.expectedProduct1.Description, suite.price, suite.expectedProduct1.Seller.ID)
	row2 := sqlmock.NewRows([]string{"id", "name", "description", "price", "seller_id"}).
		AddRow(suite.expectedProduct2.ID, suite.expectedProduct2.Name, suite.expectedProduct2.Description, suite.price, suite.expectedProduct2.Seller.ID)

	var rows = []*sqlmock.Rows{}
	rows = append(rows, row1, row2)
	prep.ExpectQuery().WillReturnRows(rows...)

	res, repoErr := suite.repo.GetAll()
	suite.NoError(repoErr)
	suite.NotNil(res)
}

func (suite *TestSuite) TestGetByID() {
	queryGetById := "SELECT id, name, description, price, seller_id FROM products WHERE id=?;"
	prep := suite.mock.ExpectPrepare(regexp.QuoteMeta(queryGetById))

	row1 := sqlmock.NewRows([]string{"id", "name", "description", "price", "seller_id"}).
		AddRow(suite.expectedProduct1.ID, suite.expectedProduct1.Name, suite.expectedProduct1.Description, suite.price, suite.expectedProduct1.Seller.ID)
	prep.ExpectQuery().WithArgs(suite.expectedProduct1.ID).WillReturnRows(row1)

	product := new(entity.Product)
	product.ID = suite.expectedProduct1.ID

	_, repoErr := suite.repo.GetByID(product)
	suite.NoError(repoErr)
	suite.NotNil(product)
}

func (suite *TestSuite) TestStore() {
	queryInsert := "INSERT INTO products(name, description, price, seller_id) VALUES(?, ?, ?, ?);"
	prep := suite.mock.ExpectPrepare(regexp.QuoteMeta(queryInsert))

	prep.ExpectExec().
		WithArgs(suite.expectedProduct1.Name, suite.expectedProduct1.Description, suite.expectedProduct1.Price, suite.expectedProduct1.Seller.ID).
		WillReturnResult(sqlmock.NewResult(suite.expectedProduct1.ID, 1))

	product := new(entity.Product)
	product.ID = suite.expectedProduct1.ID
	product.Name = suite.expectedProduct1.Name
	product.Description = suite.expectedProduct1.Description
	product.Price = suite.expectedProduct1.Price
	product.Seller = suite.expectedProduct1.Seller

	repoErr := suite.repo.Store(product)

	suite.NoError(repoErr)
	suite.NotNil(product)
}

func (suite *TestSuite) TestUpdate() {
	queryUpdate := "UPDATE products SET name=?, description=?, price=?, seller_id=? WHERE id=?;"
	prep := suite.mock.ExpectPrepare(regexp.QuoteMeta(queryUpdate))

	prep.ExpectExec().
		WithArgs(suite.expectedProduct1.Name, suite.expectedProduct1.Description, suite.expectedProduct1.Price, suite.expectedProduct1.Seller.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	product := new(entity.Product)
	product.ID = suite.expectedProduct1.ID
	product.Name = suite.expectedProduct1.Name
	product.Description = suite.expectedProduct1.Description
	product.Price = suite.expectedProduct1.Price
	product.Seller = suite.expectedProduct1.Seller

	repoErr := suite.repo.Update(product)
	suite.NoError(repoErr)
	suite.NotNil(product)
}

func (suite *TestSuite) TestDelete() {
	queryDelete := "DELETE FROM products WHERE id=?;"
	prep := suite.mock.ExpectPrepare(regexp.QuoteMeta(queryDelete))

	prep.ExpectExec().
		WithArgs(suite.expectedProduct1.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	product := new(entity.Product)
	product.ID = suite.expectedProduct1.ID

	repoErr := suite.repo.Delete(product)
	suite.NoError(repoErr)
	suite.NotNil(product)
}
