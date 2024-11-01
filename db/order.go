package db

import (
	"fmt"

	"github.com/orders-app/models"
)

type OrderDB struct {
	orders map[string]models.Order
}

// NewOrderDBService creates new order db service
func NewOrderDBService() *OrderDB {
	return &OrderDB{
		orders: make(map[string]models.Order),
	}
}

// Find order for a given order id
func (o *OrderDB) Find(id string) (models.Order, error) {
	order, ok := o.orders[id]
	if !ok {
		return models.Order{}, fmt.Errorf("no order found for %s order id", id)
	}
	return order, nil
}

// Upsert creates or updates an order in the orders database
func (o *OrderDB) Upsert(order models.Order) {
	o.orders[order.ID] = order
}
