package repo

import (
	"fmt"
	"math"

	"github.com/orders-app/db"
	"github.com/orders-app/models"
)

// repo holds all the dependencies required for repo operations
type repo struct {
	products *db.ProductDB
	orders   *db.OrderDB
	incoming chan models.Order
}

// Repo is the interface we expose to outside packages
type Repo interface {
	CreateOrder(item models.Item) (*models.Order, error)
	GetAllProducts() []models.Product
	GetOrder(id string) (models.Order, error)
}

// New creates a new Order repo with the correct database dependencies
func New() (Repo, error) {
	o := repo{
		products: db.NewProductDBService(),
		orders:   db.NewOrderDBService(),
		incoming: make(chan models.Order),
	}
	go o.processOrders()
	return &o, nil
}

// GetAllProducts returns all products in the system
func (r *repo) GetAllProducts() []models.Product {
	return r.products.GetAll()
}

// GetProduct returns the given order if one exists
func (r *repo) GetOrder(id string) (models.Order, error) {
	return r.orders.Find(id)
}

// CreateOrder creates a new order for the given item
func (r *repo) CreateOrder(item models.Item) (*models.Order, error) {
	if err := r.validateItem(item); err != nil {
		return nil, err
	}
	order := models.NewOrder(item)
	r.orders.Upsert(order)
	r.incoming <- order
	return &order, nil
}

// validateItem runs validations on a given order
func (r *repo) validateItem(item models.Item) error {
	if item.Amount < 1 {
		return fmt.Errorf("order amount must be at least 1:got %d", item.Amount)
	}
	if err := r.products.Exists(item.ProductID); err != nil {
		return fmt.Errorf("product %s does not exist", item.ProductID)
	}
	return nil
}

func (r *repo) processOrders() {
	fmt.Println("Order processing started!")

	for order := range r.incoming {
		r.processOrder(&order)
		r.orders.Upsert(order)
		fmt.Printf("Processing order %s completed\n", order.ID)
	}

	fmt.Println("Order processing stopped!")
}

// processOrder is an internal method which completes or rejects an order
func (r *repo) processOrder(order *models.Order) {
	item := order.Item
	product, err := r.products.Find(item.ProductID)
	if err != nil {
		order.Status = string(models.OrderStatus_Rejected)
		order.Error = err.Error()
		return
	}
	if product.Stock < item.Amount {
		order.Status = string(models.OrderStatus_Rejected)
		order.Error = fmt.Sprintf("not enough stock for product %s:got %d, want %d", item.ProductID, product.Stock, item.Amount)
		return
	}
	remainingStock := product.Stock - item.Amount
	product.Stock = remainingStock
	r.products.Upsert(product)

	total := math.Round(float64(order.Item.Amount)*product.Price*100) / 100
	order.Total = total
	order.Complete()
}
