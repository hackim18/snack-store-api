package usecase

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"snack-store-api/internal/cache"
	"snack-store-api/internal/constants"
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

type TransactionUseCase struct {
	DB                    *gorm.DB
	Log                   *logrus.Logger
	CustomerRepository    *repository.CustomerRepository
	ProductRepository     *repository.ProductRepository
	TransactionRepository *repository.TransactionRepository
	Cache                 cache.Cache
}

func NewTransactionUseCase(
	db *gorm.DB,
	logger *logrus.Logger,
	customerRepository *repository.CustomerRepository,
	productRepository *repository.ProductRepository,
	transactionRepository *repository.TransactionRepository,
	cacheStore cache.Cache,
) *TransactionUseCase {
	return &TransactionUseCase{
		DB:                    db,
		Log:                   logger,
		CustomerRepository:    customerRepository,
		ProductRepository:     productRepository,
		TransactionRepository: transactionRepository,
		Cache:                 cacheStore,
	}
}

func (c *TransactionUseCase) Create(
	ctx context.Context,
	request *model.CreateTransactionRequest,
) (*model.TransactionResponse, error) {
	productID, err := uuid.Parse(strings.TrimSpace(request.ProductID))
	if err != nil {
		c.Log.Warnf("Invalid product_id : %+v", err)
		return nil, utils.Error(messages.ErrInvalidIDFormat, http.StatusBadRequest, err)
	}

	transactionAt, err := time.Parse(constants.DateTimeLayout, strings.TrimSpace(request.TransactionAt))
	if err != nil {
		c.Log.Warnf("Invalid transaction_at format : %+v", err)
		return nil, utils.Error(messages.FailedInputFormat, http.StatusBadRequest, err)
	}

	customerName := strings.TrimSpace(request.CustomerName)
	if customerName == "" {
		return nil, utils.Error(messages.FailedValidationOccurred, http.StatusBadRequest, nil)
	}

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

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

	if product.StockQty < request.Qty {
		return nil, utils.Error(messages.ErrInsufficientStock, http.StatusConflict, nil)
	}

	customer, err := c.findOrCreateCustomer(tx, customerName)
	if err != nil {
		return nil, err
	}

	unitPrice := product.Price
	totalPrice := unitPrice * request.Qty
	pointsEarned := entity.PointsEarned(totalPrice)

	product.StockQty -= request.Qty
	customer.Points += pointsEarned

	if err := c.ProductRepository.Update(tx, &product); err != nil {
		c.Log.Warnf("Failed to update product stock : %+v", err)
		return nil, utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}

	if err := c.CustomerRepository.Update(tx, &customer); err != nil {
		c.Log.Warnf("Failed to update customer points : %+v", err)
		return nil, utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}

	transaction := entity.Transaction{
		CustomerID:    customer.ID,
		ProductID:     product.ID,
		Qty:           request.Qty,
		UnitPrice:     unitPrice,
		TotalPrice:    totalPrice,
		PointsEarned:  pointsEarned,
		TransactionAt: transactionAt,
	}

	if err := c.TransactionRepository.Create(tx, &transaction); err != nil {
		c.Log.Warnf("Failed to create transaction : %+v", err)
		return nil, utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}

	transaction.Customer = customer
	transaction.Product = product

	c.invalidateCaches(ctx, &product)

	return converter.TransactionToResponse(&transaction), nil
}

func (c *TransactionUseCase) List(
	ctx context.Context,
	request *model.GetTransactionRequest,
) ([]*model.TransactionResponse, model.PageMetadata, error) {
	startDate, err := time.Parse(constants.DateLayout, strings.TrimSpace(request.Start))
	if err != nil {
		c.Log.Warnf("Invalid start date : %+v", err)
		return nil, model.PageMetadata{}, utils.Error(messages.FailedInputFormat, http.StatusBadRequest, err)
	}

	endDate, err := time.Parse(constants.DateLayout, strings.TrimSpace(request.End))
	if err != nil {
		c.Log.Warnf("Invalid end date : %+v", err)
		return nil, model.PageMetadata{}, utils.Error(messages.FailedInputFormat, http.StatusBadRequest, err)
	}

	if endDate.Before(startDate) {
		return nil, model.PageMetadata{}, utils.Error(messages.InvalidRequestData, http.StatusBadRequest, nil)
	}

	endDate = endDate.AddDate(0, 0, 1)

	db := c.DB.WithContext(ctx)

	totalItem, err := c.TransactionRepository.CountByDateRange(db, startDate, endDate)
	if err != nil {
		c.Log.Warnf("Failed to count transactions : %+v", err)
		return nil, model.PageMetadata{}, utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}

	offset := (request.Page - 1) * request.PageSize
	transactions, err := c.TransactionRepository.FindByDateRange(
		db,
		startDate,
		endDate,
		request.PageSize,
		offset,
	)
	if err != nil {
		c.Log.Warnf("Failed to query transactions : %+v", err)
		return nil, model.PageMetadata{}, utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}

	responses := make([]*model.TransactionResponse, 0, len(transactions))
	for i := range transactions {
		responses = append(responses, converter.TransactionToResponse(&transactions[i]))
	}

	paging := utils.BuildPageMetadata(request.Page, request.PageSize, totalItem)
	return responses, paging, nil
}

func (c *TransactionUseCase) invalidateCaches(ctx context.Context, product *entity.Product) {
	if c.Cache == nil || product == nil {
		return
	}

	cacheKey := constants.ProductCacheKeyPrefix + product.ManufacturedDate.Format(constants.DateLayout)
	if err := c.Cache.Del(ctx, cacheKey); err != nil {
		c.Log.Warnf("Failed to invalidate product cache : %+v", err)
	}

	if err := c.Cache.DelByPrefix(ctx, constants.ReportCacheKeyPrefix); err != nil {
		c.Log.Warnf("Failed to invalidate report cache : %+v", err)
	}
}

func (c *TransactionUseCase) findOrCreateCustomer(tx *gorm.DB, name string) (entity.Customer, error) {
	var customer entity.Customer
	lowerName := strings.ToLower(strings.TrimSpace(name))

	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("lower(name) = ?", lowerName).
		Take(&customer).Error; err == nil {
		return customer, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.Log.Warnf("Failed to find customer : %+v", err)
		return entity.Customer{}, utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}

	customer = entity.Customer{
		Name:   name,
		Points: 0,
	}

	if err := tx.Create(&customer).Error; err == nil {
		return customer, nil
	}

	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("lower(name) = ?", lowerName).
		Take(&customer).Error; err != nil {
		c.Log.Warnf("Failed to refetch customer : %+v", err)
		return entity.Customer{}, utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}

	return customer, nil
}
