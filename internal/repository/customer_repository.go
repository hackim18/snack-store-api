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

func (r *CustomerRepository) FindAll(db *gorm.DB) ([]entity.Customer, error) {
	var customers []entity.Customer
	err := db.Order("created_at desc").Find(&customers).Error
	return customers, err
}
