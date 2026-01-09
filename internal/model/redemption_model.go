package model

import "github.com/google/uuid"

type CreateRedemptionRequest struct {
	CustomerName string `json:"customer_name" validate:"required"`
	ProductID    string `json:"product_id" validate:"required"`
	Qty          int    `json:"qty" validate:"required,gt=0"`
	RedeemAt     string `json:"redeem_at" validate:"required"`
}

type RedemptionResponse struct {
	ID           *uuid.UUID `json:"redemption_id,omitempty"`
	CustomerName string     `json:"customer_name,omitempty"`
	ProductName  string     `json:"product_name,omitempty"`
	Size         string     `json:"size,omitempty"`
	Qty          int        `json:"qty,omitempty"`
	PointsSpent  int        `json:"points_spent,omitempty"`
	RedeemAt     string     `json:"redeem_at,omitempty"`
}
