package repository

import (
	"snack-store-api/internal/entity"

	"github.com/sirupsen/logrus"
)

type RedemptionRepository struct {
	Repository[entity.Redemption]
	Log *logrus.Logger
}

func NewRedemptionRepository(log *logrus.Logger) *RedemptionRepository {
	return &RedemptionRepository{
		Log: log,
	}
}
