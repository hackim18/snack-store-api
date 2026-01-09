package repository

import (
	"time"

	"snack-store-api/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	Repository[entity.Transaction]
	Log *logrus.Logger
}

func NewTransactionRepository(log *logrus.Logger) *TransactionRepository {
	return &TransactionRepository{
		Log: log,
	}
}

func (r *TransactionRepository) FindByDateRange(
	db *gorm.DB,
	startDate time.Time,
	endDate time.Time,
) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	err := db.Preload("Customer").
		Preload("Product").
		Where("transaction_at >= ? AND transaction_at < ?", startDate, endDate).
		Order("transaction_at desc").
		Find(&transactions).Error
	return transactions, err
}
