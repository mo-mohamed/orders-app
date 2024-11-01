package db

import "github.com/orders-app/models"

type ProducdDB struct {
	products map[string]models.Product
}

// NewProductDBService creates a new empty products service
func NewProductDBService() *ProducdDB {
	return &ProducdDB{
		products: make(map[string]models.Product),
	}
}

// Exists checks whether a product with a given id exists
func (p *ProducdDB) Exists(id string) bool {
	_, exists := p.products[id]
	return exists
}

// Find returns a product if exists
func (p *ProducdDB) Find(id string) *models.Product {
	product, ok := p.products[id]
	if !ok {
		return nil
	}
	return &product
}

// Upsert inserts or updates a product in the database
func (p *ProducdDB) Upsert(product models.Product) {
	p.products[product.ID] = product
}

// GetAll lists all products in the database
func (p *ProducdDB) GetAll() []models.Product {
	var allProducts []models.Product

	for _, product := range p.products {
		allProducts = append(allProducts, product)
	}
	return allProducts
}
