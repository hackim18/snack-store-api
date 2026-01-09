package converter

import (
	"snack-store-api/internal/constants"
	"snack-store-api/internal/entity"
	"snack-store-api/internal/model"
)

func ProductToResponse(product *entity.Product) *model.ProductResponse {
	id := product.ID
	return &model.ProductResponse{
		ID:               &id,
		Name:             product.Name,
		Type:             product.Type,
		Flavor:           product.Flavor,
		Size:             product.Size,
		Price:            product.Price,
		StockQty:         product.StockQty,
		ManufacturedDate: product.ManufacturedDate.Format(constants.DateLayout),
	}
}
