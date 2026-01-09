package usecase

import (
	"context"
	"net/http"

	"snack-store-api/internal/messages"
	"snack-store-api/internal/model"
	"snack-store-api/internal/model/converter"
	"snack-store-api/internal/repository"
	"snack-store-api/internal/utils"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CustomerUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	CustomerRepository *repository.CustomerRepository
}

func NewCustomerUseCase(
	db *gorm.DB,
	logger *logrus.Logger,
	customerRepository *repository.CustomerRepository,
) *CustomerUseCase {
	return &CustomerUseCase{
		DB:                 db,
		Log:                logger,
		CustomerRepository: customerRepository,
	}
}

func (c *CustomerUseCase) List(ctx context.Context) ([]*model.CustomerResponse, error) {
	customers, err := c.CustomerRepository.FindAll(c.DB.WithContext(ctx))
	if err != nil {
		c.Log.Warnf("Failed to query customers : %+v", err)
		return nil, utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}

	responses := make([]*model.CustomerResponse, 0, len(customers))
	for i := range customers {
		responses = append(responses, converter.CustomerToResponse(&customers[i]))
	}

	return responses, nil
}
