package db

import (
	"fmt"
	"sync"

	"github.com/orders-app/models"
)

type OrderDB struct {
	orders sync.Map
}

// NewOrderDBService creates new order db service
func NewOrderDBService() *OrderDB {
	return &OrderDB{}
}

// Find order for a given order id
func (o *OrderDB) Find(id string) (models.Order, error) {
	order, ok := o.orders.Load(id)
	if !ok {
		return models.Order{}, fmt.Errorf("no order found for %s order id", id)
	}
	return toOrder(order), nil
}

// Upsert creates or updates an order in the orders database
func (o *OrderDB) Upsert(order models.Order) {
	o.orders.Store(order.ID, order)
}

func toOrder(o any) models.Order {
	order, ok := o.(models.Order)
	if !ok {
		panic(fmt.Errorf("error casting %v to order", o))
	}
	return order
}
