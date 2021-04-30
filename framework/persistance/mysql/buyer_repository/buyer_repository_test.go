package buyerrepo_test

import (
	"database/sql"
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
	expectedBuyer1 entity.Buyer
	expectedBuyer2 entity.Buyer
	expectedBuyer3 entity.Buyer
}

// before each test
func (suite *TestSuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New()
	suite.NoError(err)

	// encrypt password
	suite.password = "12345"
	suite.hashedPassword, err = bcrypt.GenerateFromPassword([]byte(suite.password), bcrypt.DefaultCost)
	suite.NoError(err)

	suite.expectedBuyer1 = entity.Buyer{
		ID:             1,
		Email:          "buyer1@mail.com",
		Name:           "buyer",
		Password:       string(suite.hashedPassword),
		SendingAddress: "buyer address",
	}

	suite.expectedBuyer2 = entity.Buyer{
		ID:             2,
		Email:          "buyer2@mail.com",
		Name:           "buyer",
		Password:       string(suite.hashedPassword),
		SendingAddress: "buyer address",
	}

	suite.expectedBuyer3 = entity.Buyer{
		ID:             3,
		Email:          "buyer3@mail.com",
		Name:           "buyer",
		Password:       string(suite.hashedPassword),
		SendingAddress: "buyer address",
	}
}

func TestBuyerRepo(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) TestGetAll() {
	queryGetAll := "SELECT id, email, name, sending_address FROM buyers;"
	prep := suite.mock.ExpectQuery(queryGetAll)

	row1 := sqlmock.NewRows([]string{"id", "email", "name", "sending_address"}).
		AddRow(suite.expectedBuyer1.ID, suite.expectedBuyer1.Email, suite.expectedBuyer1.Name, suite.expectedBuyer1.SendingAddress)
	row2 := sqlmock.NewRows([]string{"id", "email", "name", "sending_address"}).
		AddRow(suite.expectedBuyer2.ID, suite.expectedBuyer2.Email, suite.expectedBuyer2.Name, suite.expectedBuyer2.SendingAddress)
	row3 := sqlmock.NewRows([]string{"id", "email", "name", "sending_address"}).
		AddRow(suite.expectedBuyer3.ID, suite.expectedBuyer3.Email, suite.expectedBuyer3.Name, suite.expectedBuyer3.SendingAddress)

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

	row1 := sqlmock.NewRows([]string{"id", "email", "name", "sending_address"}).
		AddRow(suite.expectedBuyer1.ID, suite.expectedBuyer1.Email, suite.expectedBuyer1.Name, suite.expectedBuyer1.SendingAddress)
	prep.ExpectQuery().WithArgs(suite.expectedBuyer1.ID).WillReturnRows(row1)

	buyer := new(entity.Buyer)
	buyer.ID = suite.expectedBuyer1.ID

	repo := buyerrepo.NewMysqlBuyerRepository(suite.db)
	repoErr := repo.GetByID(buyer)

	suite.NoError(repoErr)
	suite.NotNil(buyer)
}

func (suite *TestSuite) TestStore() {
	queryInsert := "INSERT INTO buyers(email, name, password, sending_address) VALUES(?, ?, ?, ?);"
	prep := suite.mock.ExpectPrepare(regexp.QuoteMeta(queryInsert))

	prep.ExpectExec().
		WithArgs(suite.expectedBuyer1.Email, suite.expectedBuyer1.Name, suite.expectedBuyer1.Password, suite.expectedBuyer1.SendingAddress).
		WillReturnResult(sqlmock.NewResult(suite.expectedBuyer1.ID, 1))

	buyer := new(entity.Buyer)
	buyer.ID = suite.expectedBuyer1.ID
	buyer.Email = suite.expectedBuyer1.Email
	buyer.Name = suite.expectedBuyer1.Name
	buyer.Password = suite.expectedBuyer1.Password
	buyer.SendingAddress = suite.expectedBuyer1.SendingAddress

	repo := buyerrepo.NewMysqlBuyerRepository(suite.db)
	repoErr := repo.Store(buyer)

	suite.NoError(repoErr)
	suite.NotNil(buyer)
}

func (suite *TestSuite) TestUpdate() {
	queryUpdate := "UPDATE buyers SET email=?, name=?, sending_address=? WHERE id=?;"
	prep := suite.mock.ExpectPrepare(regexp.QuoteMeta(queryUpdate))

	prep.ExpectExec().
		WithArgs(suite.expectedBuyer1.Email, suite.expectedBuyer1.Name, suite.expectedBuyer1.SendingAddress, suite.expectedBuyer1.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	buyer := new(entity.Buyer)
	buyer.ID = suite.expectedBuyer1.ID
	buyer.Email = suite.expectedBuyer1.Email
	buyer.Name = suite.expectedBuyer1.Name
	buyer.Password = suite.expectedBuyer1.Password
	buyer.SendingAddress = suite.expectedBuyer1.SendingAddress

	repo := buyerrepo.NewMysqlBuyerRepository(suite.db)
	repoErr := repo.Update(buyer)

	suite.NoError(repoErr)
	suite.NotNil(buyer)
}

func (suite *TestSuite) TestDelete() {
	queryDelete := "DELETE FROM buyers WHERE id=?;"
	prep := suite.mock.ExpectPrepare(regexp.QuoteMeta(queryDelete))

	prep.ExpectExec().
		WithArgs(suite.expectedBuyer1.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	buyer := new(entity.Buyer)
	buyer.ID = suite.expectedBuyer1.ID

	repo := buyerrepo.NewMysqlBuyerRepository(suite.db)
	repoErr := repo.Delete(buyer)

	suite.NoError(repoErr)
	suite.NotNil(buyer)
}

func (suite *TestSuite) TestGetByEmail() {
	queryFindByEmail := "SELECT id, email, name, password, sending_address FROM buyers WHERE email=?;"
	prep := suite.mock.ExpectPrepare(regexp.QuoteMeta(queryFindByEmail))

	row1 := sqlmock.NewRows([]string{"id", "email", "name", "password", "sending_address"}).
		AddRow(suite.expectedBuyer1.ID, suite.expectedBuyer1.Email, suite.expectedBuyer1.Name, suite.expectedBuyer1.Password, suite.expectedBuyer1.SendingAddress)
	prep.ExpectQuery().WithArgs(suite.expectedBuyer1.Email).WillReturnRows(row1)

	buyer := new(entity.Buyer)
	buyer.Email = suite.expectedBuyer1.Email

	repo := buyerrepo.NewMysqlBuyerRepository(suite.db)
	repoRes, repoErr := repo.GetByEmail(buyer)

	suite.NoError(repoErr)
	suite.NotNil(repoRes)
}
