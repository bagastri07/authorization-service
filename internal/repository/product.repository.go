package repository

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/bagastri07/boilerplate-service/internal/model"
	"github.com/bagastri07/boilerplate-service/utils"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type productRepository struct {
	DB          *gorm.DB
	redisClient *redis.Client
}

func NewProductRepository(DB *gorm.DB, redis *redis.Client) model.ProductRepository {
	return &productRepository{
		DB:          DB,
		redisClient: redis,
	}
}

func (r *productRepository) Create(ctx context.Context, product *model.Product) error {
	err := r.DB.WithContext(ctx).Save(product).Error
	if err != nil {
		logrus.WithContext(ctx).WithField("product", utils.Dump(product)).Error(err)
		return err
	}

	_ = DeleteByKeys(ctx, r.redisClient, []string{model.NewProductCacheKeyFromID(product.ID)})

	return nil
}

func (r *productRepository) FindByID(ctx context.Context, ID int) (*model.Product, error) {
	logger := logrus.WithContext(ctx)

	product := new(model.Product)

	cacheKey := model.NewProductCacheKeyFromID(ID)
	cachedData, err := Get(ctx, r.redisClient, cacheKey)
	if err != nil {
		logger.WithField("cacheKey", cacheKey).Error(err.Error())
	}

	err = json.Unmarshal(cachedData, &product)
	if err == nil {
		return product, nil
	}

	product = new(model.Product)

	err = r.DB.WithContext(ctx).
		Where("id = ?", ID).
		First(product).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = SetWithExpiry(ctx, r.redisClient, cacheKey, nil)
			if err != nil {
				logger.Error(err.Error())
			}
			return nil, nil
		}
		logger.WithField("id", ID).Error(err)
		return nil, err
	}

	err = SetWithExpiry(ctx, r.redisClient, cacheKey, product)
	if err != nil {
		logger.Error(err.Error())
	}

	return product, nil
}
