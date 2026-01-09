package repository

import (
	"time"

	"snack-store-api/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ProductRepository struct {
	Repository[entity.Product]
	Log *logrus.Logger
}

func NewProductRepository(log *logrus.Logger) *ProductRepository {
	return &ProductRepository{
		Log: log,
	}
}

func (r *ProductRepository) FindByManufacturedDate(
	db *gorm.DB,
	manufacturedDate time.Time,
) ([]entity.Product, error) {
	var products []entity.Product
	err := db.Where("manufactured_date = ?", manufacturedDate).
		Order("created_at desc").
		Find(&products).Error
	return products, err
}
