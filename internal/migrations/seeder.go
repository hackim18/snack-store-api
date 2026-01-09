package migrations

import (
	"encoding/json"
	"os"
	"snack-store-api/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func Seeder(db *gorm.DB, logger *logrus.Logger) error {
	logger.Info("Seeding database...")

	seedFromJSON("internal/migrations/json/customers.json", &[]entity.Customer{}, db, logger)
	seedFromJSON("internal/migrations/json/products.json", &[]entity.Product{}, db, logger)
	seedFromJSON("internal/migrations/json/transactions.json", &[]entity.Transaction{}, db, logger)
	seedFromJSON("internal/migrations/json/redemptions.json", &[]entity.Redemption{}, db, logger)

	return nil
}

func seedFromJSON[T any](filePath string, out *[]T, db *gorm.DB, log *logrus.Logger) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Warnf("Seed file not found: %s", filePath)
		return
	}

	if err := json.Unmarshal(data, out); err != nil {
		log.Warnf("Failed to parse JSON for %s: %v", filePath, err)
		return
	}

	var count int64
	if err := db.Model(out).Count(&count).Error; err != nil {
		log.Warnf("Failed to count records for %s: %v", filePath, err)
		return
	}

	if count == 0 {
		createDB := db
		if _, ok := any(out).(*[]entity.Transaction); ok {
			createDB = createDB.Omit("Customer", "Product")
		} else if _, ok := any(out).(*[]entity.Redemption); ok {
			createDB = createDB.Omit("Customer", "Product")
		}

		if err := createDB.Create(out).Error; err != nil {
			log.Warnf("Insert failed for %s: %v", filePath, err)
		} else {
			log.Infof("Inserted seed data from %s", filePath)
		}
	} else {
		log.Infof("Skipping insert for %s: table not empty", filePath)
	}
}
