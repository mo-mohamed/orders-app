package db

import (
	"fmt"

	"github.com/orders-app/models"
	"github.com/orders-app/utils"
)

type ProductDB struct {
	products map[string]models.Product
}

// NewProductDBService creates a new empty products service
func NewProductDBService() *ProductDB {

	p := &ProductDB{
		products: make(map[string]models.Product),
	}

	utils.ImportProducts(p.products)
	return p
}

// Exists checks whether a product with a given id exists
func (p *ProductDB) Exists(id string) error {
	if _, ok := p.products[id]; !ok {
		return fmt.Errorf("no product found for id %s", id)
	}

	return nil
}

// Find returns a product if exists
func (p *ProductDB) Find(id string) (models.Product, error) {
	prod, ok := p.products[id]
	if !ok {
		return models.Product{}, fmt.Errorf("no product found for id %s", id)
	}

	return prod, nil
}

// Upsert inserts or updates a product in the database
func (p *ProductDB) Upsert(product models.Product) {
	p.products[product.ID] = product
}

// GetAll lists all products in the database
func (p *ProductDB) GetAll() []models.Product {
	var allProducts []models.Product

	for _, product := range p.products {
		allProducts = append(allProducts, product)
	}
	return allProducts
}
