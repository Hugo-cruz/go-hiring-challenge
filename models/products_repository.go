package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

// GetAllProducts returns all products with variants preloaded.
func (r *ProductsRepository) GetAllProducts() ([]Product, error) {
	var products []Product
	if err := r.db.Preload("Variants").Preload("Category").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// GetProductsByFilter returns products with pagination and optional filters.
// offset: number of products to skip (default 0)
// limit: maximum number of products to return (default 10, max 100)
// categoryID: optional category filter
// priceLessThan: optional price filter
func (r *ProductsRepository) GetProductsByFilter(offset, limit int, categoryID *uint, priceLessThan *decimal.Decimal) ([]Product, int64, error) {
	// Validate and normalize limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	var products []Product
	var total int64

	query := r.db.Preload("Variants").Preload("Category")

	// Apply filters
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}
	if priceLessThan != nil {
		query = query.Where("price < ?", priceLessThan)
	}

	// Get total count
	if err := query.Model(&Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := query.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// GetProductByCode returns a single product by its code with variants preloaded.
func (r *ProductsRepository) GetProductByCode(code string) (*Product, error) {
	var product Product
	if err := r.db.Preload("Variants").Preload("Category").Where("code = ?", code).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}
