package model

import "github.com/google/uuid"

type ReportTransactionsRequest struct {
	Start string `json:"-" validate:"required,datetime=2006-01-02"`
	End   string `json:"-" validate:"required,datetime=2006-01-02"`
}

type ReportBestSeller struct {
	ProductName string `json:"product_name,omitempty"`
	Size        string `json:"size,omitempty"`
	Flavor      string `json:"flavor,omitempty"`
	TotalQty    int    `json:"total_qty,omitempty"`
}

type ReportTransactionItem struct {
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
	IsNewCustomer bool       `json:"is_new_customer"`
}

type ReportTransactionsResponse struct {
	TotalCustomer     int64                    `json:"total_customer"`
	HasNewCustomer    bool                     `json:"has_new_customer"`
	TotalIncome       int                      `json:"total_income"`
	BestSeller        *ReportBestSeller        `json:"best_seller,omitempty"`
	TotalProductsSold int                      `json:"total_products_sold"`
	LastTransactions  []*ReportTransactionItem `json:"last_transactions,omitempty"`
}
