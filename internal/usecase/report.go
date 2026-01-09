package usecase

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"snack-store-api/internal/cache"
	"snack-store-api/internal/constants"
	"snack-store-api/internal/entity"
	"snack-store-api/internal/messages"
	"snack-store-api/internal/model"
	"snack-store-api/internal/repository"
	"snack-store-api/internal/utils"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const reportLastTransactionLim = 10

type ReportUseCase struct {
	DB               *gorm.DB
	Log              *logrus.Logger
	ReportRepository *repository.ReportRepository
	Cache            cache.Cache
}

func NewReportUseCase(
	db *gorm.DB,
	logger *logrus.Logger,
	reportRepository *repository.ReportRepository,
	cacheStore cache.Cache,
) *ReportUseCase {
	return &ReportUseCase{
		DB:               db,
		Log:              logger,
		ReportRepository: reportRepository,
		Cache:            cacheStore,
	}
}

func (c *ReportUseCase) Transactions(
	ctx context.Context,
	request *model.ReportTransactionsRequest,
) (*model.ReportTransactionsResponse, error) {
	startStr := strings.TrimSpace(request.Start)
	endStr := strings.TrimSpace(request.End)

	startDate, err := time.Parse(constants.DateLayout, startStr)
	if err != nil {
		c.Log.Warnf("Invalid start date : %+v", err)
		return nil, utils.Error(messages.FailedInputFormat, http.StatusBadRequest, err)
	}

	endDate, err := time.Parse(constants.DateLayout, endStr)
	if err != nil {
		c.Log.Warnf("Invalid end date : %+v", err)
		return nil, utils.Error(messages.FailedInputFormat, http.StatusBadRequest, err)
	}

	if endDate.Before(startDate) {
		return nil, utils.Error(messages.InvalidRequestData, http.StatusBadRequest, nil)
	}

	cacheKey := reportCacheKey(startStr, endStr)
	if c.Cache != nil {
		cached, ok, err := c.Cache.Get(ctx, cacheKey)
		if err != nil {
			c.Log.Warnf("Failed to get report cache : %+v", err)
		}
		if ok {
			var cachedResponse model.ReportTransactionsResponse
			if err := json.Unmarshal([]byte(cached), &cachedResponse); err == nil {
				return &cachedResponse, nil
			}
			c.Log.Warnf("Failed to decode report cache")
		}
	}

	endDate = endDate.AddDate(0, 0, 1)

	totalCustomer, err := c.ReportRepository.GetTotalCustomer(c.DB.WithContext(ctx), startDate, endDate)
	if err != nil {
		c.Log.Warnf("Failed to get total customer : %+v", err)
		return nil, utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}

	hasNewCustomer, err := c.ReportRepository.HasNewCustomer(c.DB.WithContext(ctx), startDate, endDate)
	if err != nil {
		c.Log.Warnf("Failed to check new customer : %+v", err)
		return nil, utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}

	totalIncome, err := c.ReportRepository.GetTotalIncome(c.DB.WithContext(ctx), startDate, endDate)
	if err != nil {
		c.Log.Warnf("Failed to get total income : %+v", err)
		return nil, utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}

	totalProductsSold, err := c.ReportRepository.GetTotalProductsSold(c.DB.WithContext(ctx), startDate, endDate)
	if err != nil {
		c.Log.Warnf("Failed to get total products sold : %+v", err)
		return nil, utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}

	bestSeller, err := c.ReportRepository.GetBestSeller(c.DB.WithContext(ctx), startDate, endDate)
	if err != nil {
		c.Log.Warnf("Failed to get best seller : %+v", err)
		return nil, utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}

	lastTransactions, err := c.ReportRepository.GetLastTransactions(
		c.DB.WithContext(ctx),
		startDate,
		endDate,
		reportLastTransactionLim,
	)
	if err != nil {
		c.Log.Warnf("Failed to get last transactions : %+v", err)
		return nil, utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}

	items := make([]*model.ReportTransactionItem, 0, len(lastTransactions))
	for i := range lastTransactions {
		items = append(items, mapReportTransaction(&lastTransactions[i]))
	}

	response := &model.ReportTransactionsResponse{
		TotalCustomer:     totalCustomer,
		HasNewCustomer:    hasNewCustomer,
		TotalIncome:       totalIncome,
		TotalProductsSold: totalProductsSold,
		LastTransactions:  items,
	}

	if bestSeller != nil {
		response.BestSeller = &model.ReportBestSeller{
			ProductName: bestSeller.ProductName,
			Size:        bestSeller.Size,
			Flavor:      bestSeller.Flavor,
			TotalQty:    bestSeller.TotalQty,
		}
	}

	if c.Cache != nil {
		payload, err := json.Marshal(response)
		if err != nil {
			c.Log.Warnf("Failed to encode report cache : %+v", err)
		} else if err := c.Cache.Set(ctx, cacheKey, string(payload), constants.ReportCacheTTL); err != nil {
			c.Log.Warnf("Failed to set report cache : %+v", err)
		}
	}

	return response, nil
}

func mapReportTransaction(transaction *entity.Transaction) *model.ReportTransactionItem {
	id := transaction.ID
	isNewCustomer := transaction.Customer.CreatedAt.Year() == transaction.TransactionAt.Year() &&
		transaction.Customer.CreatedAt.Month() == transaction.TransactionAt.Month()

	return &model.ReportTransactionItem{
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
		IsNewCustomer: isNewCustomer,
	}
}

func reportCacheKey(startDate, endDate string) string {
	return constants.ReportCacheKeyPrefix + startDate + ":" + endDate
}
