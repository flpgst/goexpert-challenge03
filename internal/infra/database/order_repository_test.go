package database

import (
	"database/sql"
	"testing"

	"github.com/flpgst/golang-studies/55-CleanArch/internal/entity"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type OrderRepositoryTestSuite struct {
	suite.Suite
	Db *sql.DB
}

func (suite *OrderRepositoryTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	suite.NoError(err)
	db.Exec("CREATE TABLE orders (id text NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, PRIMARY KEY (id))")
	suite.Db = db
}

func (suite *OrderRepositoryTestSuite) TearDownSuite() {
	suite.Db.Close()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}

func (suite *OrderRepositoryTestSuite) TestGivenAnOrder_WhenSaveIsCalled_ShouldSaveOrder() {
	order, err := entity.NewOrder("123", 10.0, 2.0)
	suite.NoError(err)
	suite.NoError(order.CalculateFinalPrice())
	repo := NewOrderRepository(suite.Db)
	err = repo.Save(order)
	suite.NoError(err)

	var orderResult entity.Order
	err = suite.Db.QueryRow("SELECT * FROM orders WHERE id = ?", order.ID).Scan(&orderResult.ID, &orderResult.Price, &orderResult.Tax, &orderResult.FinalPrice)

	suite.NoError(err)
	suite.Equal(order.ID, orderResult.ID)
	suite.Equal(order.Price, orderResult.Price)
	suite.Equal(order.Tax, orderResult.Tax)
	suite.Equal(order.FinalPrice, orderResult.FinalPrice)
}

func (suite *OrderRepositoryTestSuite) TestListOrders_WhenListIsCalled_ShouldReturnOrders() {
	order1, err := entity.NewOrder("123", 10.0, 2.0)
	suite.NoError(err)
	order2, err := entity.NewOrder("456", 20.0, 3.0)
	suite.NoError(err)
	repo := NewOrderRepository(suite.Db)
	err = repo.Save(order1)
	suite.NoError(err)
	err = repo.Save(order2)
	suite.NoError(err)

	orders, err := repo.List()
	suite.NoError(err)
	suite.Equal(2, len(orders))
	suite.Equal(order1.ID, orders[0].ID)
	suite.Equal(order2.ID, orders[1].ID)
}
