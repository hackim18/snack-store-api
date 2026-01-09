package migrations

import (
	"snack-store-api/internal/entity"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.Customer{},
		&entity.Product{},
		&entity.Transaction{},
		&entity.Redemption{},
	)
}
