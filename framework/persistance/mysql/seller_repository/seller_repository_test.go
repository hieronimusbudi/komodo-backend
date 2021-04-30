package sellerrepo_test

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/hieronimusbudi/komodo-backend/entity"
	sellerrepo "github.com/hieronimusbudi/komodo-backend/framework/persistance/mysql/seller_repository"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type TestSuite struct {
	suite.Suite
	db              *sql.DB
	mock            sqlmock.Sqlmock
	repo            entity.SellerRepository
	password        string
	hashedPassword  []byte
	expectedSeller1 entity.Seller
	expectedSeller2 entity.Seller
	expectedSeller3 entity.Seller
}

// before each test
func (suite *TestSuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New()
	suite.NoError(err)

	// initiate repo
	suite.repo = sellerrepo.NewMysqlSellerRepository(suite.db)

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

	suite.expectedSeller2 = entity.Seller{
		ID:            2,
		Email:         "seller2@mail.com",
		Name:          "seller",
		Password:      string(suite.hashedPassword),
		PickUpAddress: "pickup address",
	}

	suite.expectedSeller3 = entity.Seller{
		ID:            3,
		Email:         "seller3@mail.com",
		Name:          "seller",
		Password:      string(suite.hashedPassword),
		PickUpAddress: "pickup address",
	}
}

func TestSellerRepo(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) TestGetAll() {
	queryGetAll := "SELECT id, email, name, pickup_address FROM sellers;"
	prep := suite.mock.ExpectQuery(queryGetAll)

	row1 := sqlmock.NewRows([]string{"id", "email", "name", "pickup_address"}).
		AddRow(suite.expectedSeller1.ID, suite.expectedSeller1.Email, suite.expectedSeller1.Name, suite.expectedSeller1.PickUpAddress)
	row2 := sqlmock.NewRows([]string{"id", "email", "name", "pickup_address"}).
		AddRow(suite.expectedSeller2.ID, suite.expectedSeller2.Email, suite.expectedSeller2.Name, suite.expectedSeller2.PickUpAddress)
	row3 := sqlmock.NewRows([]string{"id", "email", "name", "pickup_address"}).
		AddRow(suite.expectedSeller3.ID, suite.expectedSeller3.Email, suite.expectedSeller3.Name, suite.expectedSeller3.PickUpAddress)

	var rows = []*sqlmock.Rows{}
	rows = append(rows, row1, row2, row3)
	prep.WillReturnRows(rows...)

	res, repoErr := suite.repo.GetAll()
	suite.NoError(repoErr)
	suite.NotNil(res)
}

func (suite *TestSuite) TestGetByID() {
	queryGetById := "SELECT id, email, name, pickup_address FROM sellers WHERE id=?;"
	prep := suite.mock.ExpectPrepare(regexp.QuoteMeta(queryGetById))

	row1 := sqlmock.NewRows([]string{"id", "email", "name", "pickup_address"}).
		AddRow(suite.expectedSeller1.ID, suite.expectedSeller1.Email, suite.expectedSeller1.Name, suite.expectedSeller1.PickUpAddress)
	prep.ExpectQuery().WithArgs(suite.expectedSeller1.ID).WillReturnRows(row1)

	seller := new(entity.Seller)
	seller.ID = suite.expectedSeller1.ID

	repoErr := suite.repo.GetByID(seller)
	suite.NoError(repoErr)
	suite.NotNil(seller)
}

func (suite *TestSuite) TestStore() {
	queryInsert := "INSERT INTO sellers(email, name, password, pickup_address) VALUES(?, ?, ?, ?);"
	prep := suite.mock.ExpectPrepare(regexp.QuoteMeta(queryInsert))

	prep.ExpectExec().
		WithArgs(suite.expectedSeller1.Email, suite.expectedSeller1.Name, suite.expectedSeller1.Password, suite.expectedSeller1.PickUpAddress).
		WillReturnResult(sqlmock.NewResult(suite.expectedSeller1.ID, 1))

	seller := new(entity.Seller)
	seller.ID = suite.expectedSeller1.ID
	seller.Email = suite.expectedSeller1.Email
	seller.Name = suite.expectedSeller1.Name
	seller.Password = suite.expectedSeller1.Password
	seller.PickUpAddress = suite.expectedSeller1.PickUpAddress

	repoErr := suite.repo.Store(seller)

	suite.NoError(repoErr)
	suite.NotNil(seller)
}

func (suite *TestSuite) TestUpdate() {
	queryUpdate := "UPDATE sellers SET email=?, name=?, pickup_address=? WHERE id=?;"
	prep := suite.mock.ExpectPrepare(regexp.QuoteMeta(queryUpdate))

	prep.ExpectExec().
		WithArgs(suite.expectedSeller1.Email, suite.expectedSeller1.Name, suite.expectedSeller1.PickUpAddress, suite.expectedSeller1.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	seller := new(entity.Seller)
	seller.ID = suite.expectedSeller1.ID
	seller.Email = suite.expectedSeller1.Email
	seller.Name = suite.expectedSeller1.Name
	seller.Password = suite.expectedSeller1.Password
	seller.PickUpAddress = suite.expectedSeller1.PickUpAddress

	repoErr := suite.repo.Update(seller)
	suite.NoError(repoErr)
	suite.NotNil(seller)
}

func (suite *TestSuite) TestDelete() {
	queryDelete := "DELETE FROM sellers WHERE id=?;"
	prep := suite.mock.ExpectPrepare(regexp.QuoteMeta(queryDelete))

	prep.ExpectExec().
		WithArgs(suite.expectedSeller1.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	seller := new(entity.Seller)
	seller.ID = suite.expectedSeller1.ID

	repoErr := suite.repo.Delete(seller)
	suite.NoError(repoErr)
	suite.NotNil(seller)
}

func (suite *TestSuite) TestGetByEmail() {
	queryFindByEmail := "SELECT id, email, name, password, pickup_address FROM sellers WHERE email=?;"
	prep := suite.mock.ExpectPrepare(regexp.QuoteMeta(queryFindByEmail))

	row1 := sqlmock.NewRows([]string{"id", "email", "name", "password", "pickup_address"}).
		AddRow(suite.expectedSeller1.ID, suite.expectedSeller1.Email, suite.expectedSeller1.Name, suite.expectedSeller1.Password, suite.expectedSeller1.PickUpAddress)
	prep.ExpectQuery().WithArgs(suite.expectedSeller1.Email).WillReturnRows(row1)

	seller := new(entity.Seller)
	seller.Email = suite.expectedSeller1.Email

	repoRes, repoErr := suite.repo.GetByEmail(seller)
	suite.NoError(repoErr)
	suite.NotNil(repoRes)
}
