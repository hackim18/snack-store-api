package converter

import (
	"snack-store-api/internal/entity"
	"snack-store-api/internal/model"
)

func RedemptionToResponse(redemption *entity.Redemption) *model.RedemptionResponse {
	id := redemption.ID
	return &model.RedemptionResponse{
		ID:           &id,
		CustomerName: redemption.Customer.Name,
		ProductName:  redemption.Product.Name,
		Size:         redemption.Product.Size,
		Qty:          redemption.Qty,
		PointsSpent:  redemption.PointsSpent,
		RedeemAt:     redemption.RedeemAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
