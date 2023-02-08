package services

import (
	"checkout-task/logger"
	"checkout-task/models"
	db "checkout-task/models/db"
	"context"
	"fmt"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"log"
	"sync"
	"time"
)

var DbConnection *gorm.DB

func ConnectDB() {
	var err error
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		Config.DBUserName, Config.DBUserPassword, Config.DBHost, Config.DBPort, Config.DBName)

	DbConnection, err = gorm.Open("postgres", dsn)
	if err != nil {
		logger.Error("Failed to connect to the Database", zap.Error(err))
	}
	logger.Info("Connected Successfully to the Database")

	DbConnection.AutoMigrate(&db.Payment{})
	DbConnection.AutoMigrate(&db.Card{})
	DbConnection.AutoMigrate(&db.Fraud{})
	DbConnection.AutoMigrate(&db.Token{})
}

var redisDefaultClient *redis.Client
var redisDefaultOnce sync.Once

var redisCache *cache.Cache
var redisCacheOnce sync.Once

func GetRedisDefaultClient() *redis.Client {
	redisDefaultOnce.Do(func() {
		redisDefaultClient = redis.NewClient(&redis.Options{
			Addr:     Config.RedisDefaultAddr,
			Password: Config.RedisPassword,
		})
	})

	return redisDefaultClient
}

func GetRedisCache() *cache.Cache {
	redisCacheOnce.Do(func() {
		redisCache = cache.New(&cache.Options{
			Redis:      GetRedisDefaultClient(),
			LocalCache: cache.NewTinyLFU(1000, time.Minute),
		})
	})

	return redisCache
}

func CheckRedisConnection() {
	redisClient := GetRedisDefaultClient()
	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		log.Fatal("Failed to connected redis:", err.Error())
	}

	log.Println("Connected to Redis!")
}

func getPaymentCacheKey(paymentID string) string {
	return "req:cache:payment:" + paymentID
}

func SetPaymentDetails(paymentID string, paymentDetails *models.PaymentResponse) {
	if !Config.UseRedis {
		return
	}

	_ = GetRedisCache().Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   getPaymentCacheKey(paymentID),
		Value: paymentDetails,
		TTL:   time.Minute,
	})
}

func GetPaymentDetailsFromCache(paymentID string) (*models.PaymentResponse, bool, error) {
	if !Config.UseRedis {
		return &models.PaymentResponse{}, false, errors.New("no redis client, set USE_REDIS in .env")
	}

	paymentDetails := &models.PaymentResponse{}
	paymentCacheKey := getPaymentCacheKey(paymentID)
	err := GetRedisCache().Get(context.TODO(), paymentCacheKey, paymentDetails)
	if err != nil {
		return paymentDetails, false, err
	}
	return paymentDetails, true, nil
}
