package converter

import (
	"snack-store-api/internal/entity"
	"snack-store-api/internal/model"
)

func CustomerToResponse(customer *entity.Customer) *model.CustomerResponse {
	return &model.CustomerResponse{
		Name:   customer.Name,
		Points: customer.Points,
	}
}
