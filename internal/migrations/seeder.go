package migrations

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func Seeder(db *gorm.DB, logger *logrus.Logger) error {
	logger.Info("Seeding database...")
	return nil
}
