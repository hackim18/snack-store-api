package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Redemption struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CustomerID  uuid.UUID `gorm:"type:uuid;not null;index:redemptions_customer_id_idx"`
	Customer    Customer  `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	ProductID   uuid.UUID `gorm:"type:uuid;not null;index:redemptions_product_id_idx"`
	Product     Product   `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	Qty         int       `gorm:"not null;check:qty > 0"`
	PointsSpent int       `gorm:"column:points_spent;not null;check:points_spent >= 0"`
	RedeemAt    time.Time `gorm:"column:redeem_at;not null;index:redemptions_redeem_at_idx"`
	CreatedAt   time.Time `gorm:"not null;default:now()"`
}

func (r *Redemption) TableName() string {
	return "redemptions"
}

func (r *Redemption) BeforeCreate(_ *gorm.DB) (err error) {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}

	return
}
