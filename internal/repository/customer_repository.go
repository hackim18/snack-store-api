package repository

import (
	"snack-store-api/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CustomerRepository struct {
	Repository[entity.Customer]
	Log *logrus.Logger
}

func NewCustomerRepository(log *logrus.Logger) *CustomerRepository {
	return &CustomerRepository{
		Log: log,
	}
}

func (r *CustomerRepository) FindAll(db *gorm.DB, limit int, offset int) ([]entity.Customer, error) {
	var customers []entity.Customer
	err := db.Order("created_at desc").Limit(limit).Offset(offset).Find(&customers).Error
	return customers, err
}

func (r *CustomerRepository) CountAll(db *gorm.DB) (int64, error) {
	var total int64
	err := db.Model(&entity.Customer{}).Count(&total).Error
	return total, err
}
