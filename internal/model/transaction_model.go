package model

import "github.com/google/uuid"

type CreateTransactionRequest struct {
	CustomerName  string `json:"customer_name" validate:"required"`
	ProductID     string `json:"product_id" validate:"required"`
	Qty           int    `json:"qty" validate:"required,gt=0"`
	TransactionAt string `json:"transaction_at" validate:"required"`
}

type GetTransactionRequest struct {
	Start string `json:"-" validate:"required,datetime=2006-01-02"`
	End   string `json:"-" validate:"required,datetime=2006-01-02"`
}

type TransactionResponse struct {
	ID            *uuid.UUID `json:"transaction_id,omitempty"`
	CustomerName  string     `json:"customer_name,omitempty"`
	ProductName   string     `json:"product_name,omitempty"`
	Size          string     `json:"size,omitempty"`
	Flavor        string     `json:"flavor,omitempty"`
	Qty           int        `json:"qty,omitempty"`
	UnitPrice     int        `json:"unit_price,omitempty"`
	TotalPrice    int        `json:"total_price,omitempty"`
	PointsEarned  int        `json:"points_earned,omitempty"`
	TransactionAt string     `json:"transaction_at,omitempty"`
}
