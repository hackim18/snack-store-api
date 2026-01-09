package usecase

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"snack-store-api/internal/cache"
	"snack-store-api/internal/entity"
	"snack-store-api/internal/messages"
	"snack-store-api/internal/model"
	"snack-store-api/internal/model/converter"
	"snack-store-api/internal/repository"
	"snack-store-api/internal/utils"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RedemptionUseCase struct {
	DB                   *gorm.DB
	Log                  *logrus.Logger
	CustomerRepository   *repository.CustomerRepository
	ProductRepository    *repository.ProductRepository
	RedemptionRepository *repository.RedemptionRepository
	Cache                cache.Cache
}

func NewRedemptionUseCase(
	db *gorm.DB,
	logger *logrus.Logger,
	customerRepository *repository.CustomerRepository,
	productRepository *repository.ProductRepository,
	redemptionRepository *repository.RedemptionRepository,
	cacheStore cache.Cache,
) *RedemptionUseCase {
	return &RedemptionUseCase{
		DB:                   db,
		Log:                  logger,
		CustomerRepository:   customerRepository,
		ProductRepository:    productRepository,
		RedemptionRepository: redemptionRepository,
		Cache:                cacheStore,
	}
}

func (c *RedemptionUseCase) Create(
	ctx context.Context,
	request *model.CreateRedemptionRequest,
) (*model.RedemptionResponse, error) {
	productID, err := uuid.Parse(strings.TrimSpace(request.ProductID))
	if err != nil {
		c.Log.Warnf("Invalid product_id : %+v", err)
		return nil, utils.Error(messages.ErrInvalidIDFormat, http.StatusBadRequest, err)
	}

	redeemAt, err := time.Parse(time.RFC3339, strings.TrimSpace(request.RedeemAt))
	if err != nil {
		c.Log.Warnf("Invalid redeem_at format : %+v", err)
		return nil, utils.Error(messages.FailedInputFormat, http.StatusBadRequest, err)
	}

	customerName := strings.TrimSpace(request.CustomerName)
	if customerName == "" {
		return nil, utils.Error(messages.FailedValidationOccurred, http.StatusBadRequest, nil)
	}

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var customer entity.Customer
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("lower(name) = ?", strings.ToLower(customerName)).
		Take(&customer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.Error(messages.StatusNotFound, http.StatusNotFound, err)
		}
		c.Log.Warnf("Failed to lock customer : %+v", err)
		return nil, utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}

	var product entity.Product
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", productID).
		Take(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.Error(messages.StatusNotFound, http.StatusNotFound, err)
		}
		c.Log.Warnf("Failed to lock product : %+v", err)
		return nil, utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}

	pointsCost := entity.PointsCost(product.Size)
	if pointsCost == 0 {
		return nil, utils.Error(messages.InvalidRequestData, http.StatusBadRequest, nil)
	}

	totalPoints := pointsCost * request.Qty
	if customer.Points < totalPoints {
		return nil, utils.Error(messages.ErrInsufficientPoints, http.StatusConflict, nil)
	}

	if product.StockQty < request.Qty {
		return nil, utils.Error(messages.ErrInsufficientStock, http.StatusConflict, nil)
	}

	customer.Points -= totalPoints
	product.StockQty -= request.Qty

	if err := c.CustomerRepository.Update(tx, &customer); err != nil {
		c.Log.Warnf("Failed to update customer points : %+v", err)
		return nil, utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}

	if err := c.ProductRepository.Update(tx, &product); err != nil {
		c.Log.Warnf("Failed to update product stock : %+v", err)
		return nil, utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}

	redemption := entity.Redemption{
		CustomerID:  customer.ID,
		ProductID:   product.ID,
		Qty:         request.Qty,
		PointsSpent: totalPoints,
		RedeemAt:    redeemAt,
	}

	if err := c.RedemptionRepository.Create(tx, &redemption); err != nil {
		c.Log.Warnf("Failed to create redemption : %+v", err)
		return nil, utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}

	redemption.Customer = customer
	redemption.Product = product

	c.invalidateCaches(ctx, &product)

	return converter.RedemptionToResponse(&redemption), nil
}

func (c *RedemptionUseCase) invalidateCaches(ctx context.Context, product *entity.Product) {
	if c.Cache == nil || product == nil {
		return
	}

	cacheKey := "products:date:" + product.ManufacturedDate.Format("2006-01-02")
	if err := c.Cache.Del(ctx, cacheKey); err != nil {
		c.Log.Warnf("Failed to invalidate product cache : %+v", err)
	}

	if err := c.Cache.DelByPrefix(ctx, "report:transactions:"); err != nil {
		c.Log.Warnf("Failed to invalidate report cache : %+v", err)
	}
}
