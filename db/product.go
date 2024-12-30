package db

import (
	"fmt"
	"sync"

	"github.com/orders-app/models"
	"github.com/orders-app/utils"
)

type ProductDB struct {
	products sync.Map
}

// NewProductDBService creates a new empty products service
func NewProductDBService() *ProductDB {

	p := &ProductDB{}

	utils.ImportProducts(&p.products)
	return p
}

// Exists checks whether a product with a given id exists
func (p *ProductDB) Exists(id string) error {
	if _, ok := p.products.Load(id); !ok {
		return fmt.Errorf("no product found for id %s", id)
	}

	return nil
}

// Find returns a product if exists
func (p *ProductDB) Find(id string) (models.Product, error) {
	prod, ok := p.products.Load(id)
	if !ok {
		return models.Product{}, fmt.Errorf("no product found for id %s", id)
	}

	return toProduct(prod), nil
}

// Upsert inserts or updates a product in the database
func (p *ProductDB) Upsert(product models.Product) {
	p.products.Store(product.ID, product)
}

// GetAll lists all products in the database
func (p *ProductDB) GetAll() []models.Product {
	var allProducts []models.Product

	p.products.Range(func(key, value any) bool {
		allProducts = append(allProducts, toProduct(value))
		return true
	})

	return allProducts
}

func toProduct(p any) models.Product {
	product, ok := p.(models.Product)

	if !ok {
		panic(fmt.Errorf("error casting %v to product", p))
	}
	return product
}
