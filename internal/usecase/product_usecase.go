package usecase

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"snack-store-api/internal/cache"
	"snack-store-api/internal/entity"
	"snack-store-api/internal/messages"
	"snack-store-api/internal/model"
	"snack-store-api/internal/model/converter"
	"snack-store-api/internal/repository"
	"snack-store-api/internal/utils"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const productCacheTTL = 5 * time.Minute

type ProductUseCase struct {
	DB                *gorm.DB
	Log               *logrus.Logger
	ProductRepository *repository.ProductRepository
	Cache             cache.Cache
}

func NewProductUseCase(
	db *gorm.DB,
	logger *logrus.Logger,
	productRepository *repository.ProductRepository,
	cacheStore cache.Cache,
) *ProductUseCase {
	return &ProductUseCase{
		DB:                db,
		Log:               logger,
		ProductRepository: productRepository,
		Cache:             cacheStore,
	}
}

func (c *ProductUseCase) ListByDate(
	ctx context.Context,
	request *model.GetProductRequest,
) ([]*model.ProductResponse, error) {
	manufacturedDate, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		c.Log.Warnf("Invalid manufactured_date format : %+v", err)
		return nil, utils.Error(messages.FailedInputFormat, http.StatusBadRequest, err)
	}

	cacheKey := productCacheKey(request.Date)
	if c.Cache != nil {
		cached, ok, err := c.Cache.Get(ctx, cacheKey)
		if err != nil {
			c.Log.Warnf("Failed to get product cache : %+v", err)
		}

		if ok {
			var cachedResponses []*model.ProductResponse
			if err := json.Unmarshal([]byte(cached), &cachedResponses); err == nil {
				return cachedResponses, nil
			}
			c.Log.Warnf("Failed to decode product cache")
		}
	}

	products, err := c.ProductRepository.FindByManufacturedDate(c.DB.WithContext(ctx), manufacturedDate)
	if err != nil {
		c.Log.Warnf("Failed to query products : %+v", err)
		return nil, utils.Error(messages.InternalServerError, http.StatusInternalServerError, err)
	}

	responses := make([]*model.ProductResponse, 0, len(products))
	for i := range products {
		responses = append(responses, converter.ProductToResponse(&products[i]))
	}

	if c.Cache != nil {
		payload, err := json.Marshal(responses)
		if err != nil {
			c.Log.Warnf("Failed to encode product cache : %+v", err)
		} else if err := c.Cache.Set(ctx, cacheKey, string(payload), productCacheTTL); err != nil {
			c.Log.Warnf("Failed to set product cache : %+v", err)
		}
	}

	return responses, nil
}

func (c *ProductUseCase) Create(
	ctx context.Context,
	request *model.CreateProductRequest,
) (*model.ProductResponse, error) {
	manufacturedDate, err := time.Parse("2006-01-02", request.ManufacturedDate)
	if err != nil {
		c.Log.Warnf("Invalid manufactured_date format : %+v", err)
		return nil, utils.Error(messages.FailedInputFormat, http.StatusBadRequest, err)
	}

	product := entity.Product{
		Name:             request.Name,
		Type:             request.Type,
		Flavor:           request.Flavor,
		Size:             request.Size,
		Price:            request.Price,
		StockQty:         request.StockQty,
		ManufacturedDate: manufacturedDate,
	}

	if err := c.ProductRepository.Create(c.DB.WithContext(ctx), &product); err != nil {
		c.Log.Warnf("Failed to create product : %+v", err)
		return nil, utils.Error(messages.ErrCreateProduct, http.StatusInternalServerError, err)
	}

	if c.Cache != nil {
		cacheKey := productCacheKey(request.ManufacturedDate)
		if err := c.Cache.Del(ctx, cacheKey); err != nil {
			c.Log.Warnf("Failed to invalidate product cache : %+v", err)
		}
	}

	return converter.ProductToResponse(&product), nil
}

func productCacheKey(date string) string {
	return "products:date:" + date
}
