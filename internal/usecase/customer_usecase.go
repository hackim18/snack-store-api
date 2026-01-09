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

func (c *CustomerUseCase) List(
	ctx context.Context,
	request *model.GetCustomerRequest,
) ([]*model.CustomerResponse, model.PageMetadata, error) {
	db := c.DB.WithContext(ctx)

	totalItem, err := c.CustomerRepository.CountAll(db)
	if err != nil {
		c.Log.Warnf("Failed to count customers : %+v", err)
		return nil, model.PageMetadata{}, utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}

	offset := (request.Page - 1) * request.PageSize
	customers, err := c.CustomerRepository.FindAll(db, request.PageSize, offset)
	if err != nil {
		c.Log.Warnf("Failed to query customers : %+v", err)
		return nil, model.PageMetadata{}, utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}

	responses := make([]*model.CustomerResponse, 0, len(customers))
	for i := range customers {
		responses = append(responses, converter.CustomerToResponse(&customers[i]))
	}

	paging := utils.BuildPageMetadata(request.Page, request.PageSize, totalItem)
	return responses, paging, nil
}
