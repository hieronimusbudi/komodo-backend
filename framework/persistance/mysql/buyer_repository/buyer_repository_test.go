package buyerrepo_test

import (
	"database/sql"
	"log"
	"regexp"
	"testing"

	"github.com/hieronimusbudi/komodo-backend/entity"
	buyerrepo "github.com/hieronimusbudi/komodo-backend/framework/persistance/mysql/buyer_repository"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type TestSuite struct {
	suite.Suite
	db             *sql.DB
	mock           sqlmock.Sqlmock
	password       string
	hashedPassword []byte
}

// before each test
func (suite *TestSuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	// encrypt password
	suite.password = "12345"
	suite.hashedPassword, err = bcrypt.GenerateFromPassword([]byte(suite.password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func TestBuyerRepo(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) TestGetAll() {
	queryGetAll := "SELECT id, email, name, sending_address FROM buyers;"
	prep := suite.mock.ExpectQuery(queryGetAll)

	row1 := sqlmock.NewRows([]string{"id", "email", "name", "sending_address"}).
		AddRow(1, "buyer1@mail.com", "buyer", "buyer address")
	row2 := sqlmock.NewRows([]string{"id", "email", "name", "sending_address"}).
		AddRow(2, "buyer2@mail.com", "buyer2", "buyer address")
	row3 := sqlmock.NewRows([]string{"id", "email", "name", "sending_address"}).
		AddRow(3, "buyer3@mail.com", "buyer3", "buyer address")

	var rows = []*sqlmock.Rows{}
	rows = append(rows, row1, row2, row3)
	prep.WillReturnRows(rows...)

	repo := buyerrepo.NewMysqlBuyerRepository(suite.db)
	res, repoErr := repo.GetAll()

	suite.NoError(repoErr)
	suite.NotNil(res)
}

func (suite *TestSuite) TestGetByID() {
	queryGetById := "SELECT id, email, name, sending_address FROM buyers WHERE id=?;"
	prep := suite.mock.ExpectPrepare(regexp.QuoteMeta(queryGetById))
	buyerID := int64(1)

	row1 := sqlmock.NewRows([]string{"id", "email", "name", "sending_address"}).
		AddRow(buyerID, "buyer1@mail.com", "buyer", "buyer addressx")
	prep.ExpectQuery().WithArgs(buyerID).WillReturnRows(row1)

	buyer := new(entity.Buyer)
	buyer.ID = buyerID

	repo := buyerrepo.NewMysqlBuyerRepository(suite.db)
	repoErr := repo.GetByID(buyer)

	suite.NoError(repoErr)
	suite.NotNil(buyer)
}

func (suite *TestSuite) TestStore() {
	queryInsert := "INSERT INTO buyers(email, name, password, sending_address) VALUES(?, ?, ?, ?);"
	prep := suite.mock.ExpectPrepare(regexp.QuoteMeta(queryInsert))

	expectedBuyers := entity.Buyer{
		ID:             1,
		Email:          "buyer1@mail.com",
		Name:           "buyer",
		Password:       string(suite.hashedPassword),
		SendingAddress: "buyer address",
	}

	prep.ExpectExec().
		WithArgs(expectedBuyers.Email, expectedBuyers.Name, expectedBuyers.Password, expectedBuyers.SendingAddress).
		WillReturnResult(sqlmock.NewResult(expectedBuyers.ID, 1))

	buyer := new(entity.Buyer)
	buyer.ID = expectedBuyers.ID
	buyer.Email = expectedBuyers.Email
	buyer.Name = expectedBuyers.Name
	buyer.Password = expectedBuyers.Password
	buyer.SendingAddress = expectedBuyers.SendingAddress

	repo := buyerrepo.NewMysqlBuyerRepository(suite.db)
	repoErr := repo.Store(buyer)

	suite.NoError(repoErr)
	suite.NotNil(buyer)
}

func (suite *TestSuite) TestUpdate() {
	queryUpdate := "UPDATE buyers SET email=?, name=?, sending_address=? WHERE id=?;"
	prep := suite.mock.ExpectPrepare(regexp.QuoteMeta(queryUpdate))

	expectedBuyers := entity.Buyer{
		ID:             1,
		Email:          "buyer1@mail.com",
		Name:           "buyer",
		Password:       string(suite.hashedPassword),
		SendingAddress: "buyer address",
	}

	prep.ExpectExec().
		WithArgs(expectedBuyers.Email, expectedBuyers.Name, expectedBuyers.SendingAddress, expectedBuyers.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	buyer := new(entity.Buyer)
	buyer.ID = expectedBuyers.ID
	buyer.Email = expectedBuyers.Email
	buyer.Name = expectedBuyers.Name
	buyer.Password = expectedBuyers.Password
	buyer.SendingAddress = expectedBuyers.SendingAddress

	repo := buyerrepo.NewMysqlBuyerRepository(suite.db)
	repoErr := repo.Update(buyer)

	suite.NoError(repoErr)
	suite.NotNil(buyer)
}

func (suite *TestSuite) TestDelete() {
	queryDelete := "DELETE FROM buyers WHERE id=?;"
	prep := suite.mock.ExpectPrepare(regexp.QuoteMeta(queryDelete))

	expectedBuyers := entity.Buyer{
		ID:             1,
		Email:          "buyer1@mail.com",
		Name:           "buyer",
		Password:       string(suite.hashedPassword),
		SendingAddress: "buyer address",
	}

	prep.ExpectExec().
		WithArgs(expectedBuyers.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	buyer := new(entity.Buyer)
	buyer.ID = expectedBuyers.ID

	repo := buyerrepo.NewMysqlBuyerRepository(suite.db)
	repoErr := repo.Delete(buyer)

	suite.NoError(repoErr)
	suite.NotNil(buyer)
}

func (suite *TestSuite) TestGetByEmail() {
	queryFindByEmail := "SELECT id, email, name, password, sending_address FROM buyers WHERE email=?;"
	prep := suite.mock.ExpectPrepare(regexp.QuoteMeta(queryFindByEmail))

	expectedBuyers := entity.Buyer{
		ID:             1,
		Email:          "buyer1@mail.com",
		Name:           "buyer",
		Password:       string(suite.hashedPassword),
		SendingAddress: "buyer address",
	}

	row1 := sqlmock.NewRows([]string{"id", "email", "name", "password", "sending_address"}).
		AddRow(expectedBuyers.ID, expectedBuyers.Email, expectedBuyers.Name, expectedBuyers.Password, expectedBuyers.SendingAddress)
	prep.ExpectQuery().WithArgs(expectedBuyers.Email).WillReturnRows(row1)

	buyer := new(entity.Buyer)
	buyer.Email = expectedBuyers.Email

	repo := buyerrepo.NewMysqlBuyerRepository(suite.db)
	repoRes, repoErr := repo.GetByEmail(buyer)

	suite.NoError(repoErr)
	suite.NotNil(repoRes)
}
