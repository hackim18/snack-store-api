package repository

import (
	"time"

	"snack-store-api/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BestSellerRow struct {
	ProductName string `gorm:"column:product_name"`
	Size        string `gorm:"column:size"`
	Flavor      string `gorm:"column:flavor"`
	TotalQty    int    `gorm:"column:total_qty"`
}

type ReportRepository struct {
	Log *logrus.Logger
}

func NewReportRepository(log *logrus.Logger) *ReportRepository {
	return &ReportRepository{Log: log}
}

func (r *ReportRepository) GetTotalCustomer(db *gorm.DB, startDate, endDate time.Time) (int64, error) {
	var total int64
	err := db.Model(&entity.Transaction{}).
		Where("transaction_at >= ? AND transaction_at < ?", startDate, endDate).
		Distinct("customer_id").
		Count(&total).Error
	return total, err
}

func (r *ReportRepository) HasNewCustomer(db *gorm.DB, startDate, endDate time.Time) (bool, error) {
	var exists bool
	err := db.Raw(`
SELECT EXISTS (
  SELECT 1
  FROM transactions t
  JOIN customers c ON c.id = t.customer_id
  WHERE t.transaction_at >= ? AND t.transaction_at < ?
    AND date_trunc('month', c.created_at) = date_trunc('month', t.transaction_at)
)`, startDate, endDate).Scan(&exists).Error
	return exists, err
}

func (r *ReportRepository) GetTotalIncome(db *gorm.DB, startDate, endDate time.Time) (int, error) {
	var total int64
	err := db.Model(&entity.Transaction{}).
		Select("COALESCE(SUM(total_price), 0)").
		Where("transaction_at >= ? AND transaction_at < ?", startDate, endDate).
		Scan(&total).Error
	return int(total), err
}

func (r *ReportRepository) GetTotalProductsSold(db *gorm.DB, startDate, endDate time.Time) (int, error) {
	var total int64
	err := db.Model(&entity.Transaction{}).
		Select("COALESCE(SUM(qty), 0)").
		Where("transaction_at >= ? AND transaction_at < ?", startDate, endDate).
		Scan(&total).Error
	return int(total), err
}

func (r *ReportRepository) GetBestSeller(db *gorm.DB, startDate, endDate time.Time) (*BestSellerRow, error) {
	var row BestSellerRow
	err := db.Raw(`
SELECT p.name AS product_name, p.size, p.flavor, SUM(t.qty) AS total_qty
FROM transactions t
JOIN products p ON p.id = t.product_id
WHERE t.transaction_at >= ? AND t.transaction_at < ?
GROUP BY p.id, p.name, p.size, p.flavor
ORDER BY total_qty DESC
LIMIT 1
`, startDate, endDate).Scan(&row).Error
	if err != nil {
		return nil, err
	}

	if row.ProductName == "" && row.TotalQty == 0 {
		return nil, nil
	}

	return &row, nil
}

func (r *ReportRepository) GetLastTransactions(
	db *gorm.DB,
	startDate time.Time,
	endDate time.Time,
	limit int,
) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	err := db.Preload("Customer").
		Preload("Product").
		Where("transaction_at >= ? AND transaction_at < ?", startDate, endDate).
		Order("transaction_at desc").
		Limit(limit).
		Find(&transactions).Error
	return transactions, err
}
