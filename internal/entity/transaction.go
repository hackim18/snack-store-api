package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CustomerID    uuid.UUID `gorm:"type:uuid;not null;index:transactions_customer_id_idx"`
	Customer      Customer  `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	ProductID     uuid.UUID `gorm:"type:uuid;not null;index:transactions_product_id_idx;index:transactions_product_time_idx,priority:1"`
	Product       Product   `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	Qty           int       `gorm:"not null;check:qty > 0"`
	UnitPrice     int       `gorm:"column:unit_price;not null;check:unit_price >= 0"`
	TotalPrice    int       `gorm:"column:total_price;not null;check:total_price >= 0"`
	PointsEarned  int       `gorm:"column:points_earned;not null;check:points_earned >= 0"`
	TransactionAt time.Time `gorm:"column:transaction_at;not null;index:transactions_transaction_at_idx;index:transactions_product_time_idx,priority:2"`
	CreatedAt     time.Time `gorm:"not null;default:now()"`
}

func (t *Transaction) TableName() string {
	return "transactions"
}

func (t *Transaction) BeforeCreate(_ *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}

	return
}
