package repo_test

import (
	"os"
	"testing"
	"time"

	"github.com/orders-app/logger"
	"github.com/orders-app/models"
	"github.com/orders-app/repo"
	"github.com/stretchr/testify/assert"
)

const existingProduct = "MWBLU"

func TestMain(m *testing.M) {
	if err := os.Chdir(".."); err != nil {
		panic(err)
	}
	logger.InitLogger("test")
	code := m.Run()
	os.Exit(code)
}

func Test_CreateOrder(t *testing.T) {
	t.Run("create & not enough stock order", func(t *testing.T) {
		rp := initRepo(t)

		item := models.Item{
			ProductID: existingProduct,
			Amount:    500,
		}
		order, _ := rp.CreateOrder(item)
		assert.NotNil(t, order)
		assert.Equal(t, string(models.OrderStatus_new), order.Status)
		assert.Equal(t, item, order.Item)
		assert.Equal(t, "", order.Error)

		// wait for the order to be processed
		time.Sleep(time.Millisecond * 100)

		dbOrder, err := rp.GetOrder(order.ID)
		assert.Nil(t, err)
		assert.NotNil(t, order)
		assert.Equal(t, string(models.OrderStatus_Rejected), dbOrder.Status)
		assert.Equal(t, item, order.Item)
		assert.Contains(t, dbOrder.Error, "not enough stock")
	})
	t.Run("create & invalid item order", func(t *testing.T) {
		rp := initRepo(t)

		item := models.Item{
			ProductID: "blablabla",
			Amount:    5,
		}
		order, err := rp.CreateOrder(item)
		assert.Nil(t, order)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "does not exist")
	})

	t.Run("create & negative stock order", func(t *testing.T) {
		rp := initRepo(t)

		item := models.Item{
			ProductID: existingProduct,
			Amount:    -5,
		}
		order, err := rp.CreateOrder(item)
		assert.Nil(t, order)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "order amount must be at least 1")
	})

}

func Test_GetOrder(t *testing.T) {
	t.Run("existing order", func(t *testing.T) {
		rp := initRepo(t)
		item := models.Item{
			ProductID: existingProduct,
			Amount:    5,
		}

		order, err := rp.CreateOrder(item)
		assert.Nil(t, err)
		assert.NotNil(t, order)

		fetchedOrder, err := rp.GetOrder(order.ID)
		assert.Nil(t, err)
		assert.NotNil(t, order)
		assert.Equal(t, *order, fetchedOrder)
	})

	t.Run("non-existing order", func(t *testing.T) {
		rp := initRepo(t)
		item := models.Item{
			ProductID: existingProduct,
			Amount:    5,
		}

		order, err := rp.CreateOrder(item)
		assert.Nil(t, err)
		assert.NotNil(t, order)

		_, err = rp.GetOrder("blablabla")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "no order found")
	})
}

func Test_GetAllProducts(t *testing.T) {
	t.Run("get products", func(t *testing.T) {
		rp := initRepo(t)
		products := rp.GetAllProducts()
		assert.Greater(t, len(products), 0)
	})
}

func initRepo(t *testing.T) repo.Repo {
	rp, err := repo.New()
	assert.Nil(t, err)
	return rp
}
