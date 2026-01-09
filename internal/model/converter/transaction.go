package converter

import (
	"snack-store-api/internal/constants"
	"snack-store-api/internal/entity"
	"snack-store-api/internal/model"
)

func TransactionToResponse(transaction *entity.Transaction) *model.TransactionResponse {
	id := transaction.ID
	return &model.TransactionResponse{
		ID:            &id,
		CustomerName:  transaction.Customer.Name,
		ProductName:   transaction.Product.Name,
		Size:          transaction.Product.Size,
		Flavor:        transaction.Product.Flavor,
		Qty:           transaction.Qty,
		UnitPrice:     transaction.UnitPrice,
		TotalPrice:    transaction.TotalPrice,
		PointsEarned:  transaction.PointsEarned,
		TransactionAt: transaction.TransactionAt.Format(constants.DateTimeLayout),
	}
}
