package models

import (
	"gorm.io/gorm"
)

type CategoriesRepository struct {
	db *gorm.DB
}

func NewCategoriesRepository(db *gorm.DB) *CategoriesRepository {
	return &CategoriesRepository{
		db: db,
	}
}

func (r *CategoriesRepository) GetAll() ([]Category, error) {
	var categories []Category
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoriesRepository) Create(category *Category) error {
	return r.db.Create(category).Error
}

func (r *CategoriesRepository) FindByCode(code string) (*Category, error) {
	var category Category
	if err := r.db.Where("code = ?", code).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}
