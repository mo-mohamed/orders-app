package repo

import (
	"context"
	"fmt"
	"math"

	"github.com/orders-app/db"
	"github.com/orders-app/models"
	"github.com/orders-app/stats"
)

// repo holds all the dependencies required for repo operations
type repo struct {
	products *db.ProductDB
	orders   *db.OrderDB
	incoming chan models.Order
	stats    stats.StatsService
	done     chan struct{}
	isOpen   bool
}

// Repo is the interface we expose to outside packages
type Repo interface {
	CreateOrder(item models.Item) (*models.Order, error)
	GetAllProducts() []models.Product
	GetOrder(id string) (models.Order, error)
	Close()
	Open()
	IsAppOpen() bool
	GetOrderStats(ctx context.Context) (models.Statistics, error)
	RequestReversal(orderId string) (*models.Order, error)
}

// New creates a new Order repo with the correct database dependencies
func New() (Repo, error) {
	processed := make(chan models.Order, stats.WorkerCount)
	done := make(chan struct{})
	statsService := stats.New(processed, done)
	o := repo{
		products: db.NewProductDBService(),
		orders:   db.NewOrderDBService(),
		incoming: make(chan models.Order),
		done:     make(chan struct{}),
		isOpen:   true,
		stats:    statsService,
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

	select {
	case r.incoming <- order:
		r.orders.Upsert(order)
		return &order, nil
	case <-r.done:
		return nil, fmt.Errorf("orders app is closed, please try gain later")
	}
}

func (r *repo) Close() {
	close(r.done)
	r.isOpen = false
}

func (r *repo) Open() {
	r.incoming = make(chan models.Order)
	r.done = make(chan struct{})
	r.isOpen = true
	go r.processOrders()
}

// GetOrderStats returns the order statistics of the orders app
func (r repo) GetOrderStats(ctx context.Context) (models.Statistics, error) {
	select {
	case s := <-r.stats.GetStats(ctx):
		return s, nil
	case <-ctx.Done():
		return models.Statistics{}, ctx.Err()
	}
}

func (r *repo) IsAppOpen() bool { return r.isOpen }

// RequestReversal fetches an existing order and updates it for reversal
func (r repo) RequestReversal(orderId string) (*models.Order, error) {
	// try to find the order first
	order, err := r.orders.Find(orderId)
	if err != nil {
		return nil, err
	}
	if order.Status != string(models.OrderStatus_Completed) {
		return nil, fmt.Errorf("order status is %s, only completed orders can be requested for reversal", order.Status)
	}
	// set reversal requested
	order.Status = string(models.OrderStatus_ReversalRequested)
	// place the order on the incoming orders channel
	select {
	case r.incoming <- order:
		r.orders.Upsert(order)
		return &order, nil
	case <-r.done:
		return nil, fmt.Errorf("sorry, the orders app is closed")
	}
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

	for {
		select {
		case order := <-r.incoming:
			r.processOrder(&order)
			r.orders.Upsert(order)
			fmt.Printf("Processing order %s completed\n", order.ID)
		case <-r.done:
			fmt.Println("Order processing stopped!")
			return
		}
	}
}

// processOrder is an internal method which completes or rejects an order
func (r *repo) processOrder(order *models.Order) {
	fetchedOrder, err := r.orders.Find(order.ID)
	if err != nil || fetchedOrder.Status != string(models.OrderStatus_Completed) {
		fmt.Println("duplicate reversal on order ", order.ID)
	}
	item := order.Item
	if order.Status == string(models.OrderStatus_ReversalRequested) {
		item.Amount = -item.Amount
	}
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
