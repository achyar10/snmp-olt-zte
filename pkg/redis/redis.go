package redis

import (
	"os"
	"time"

	"github.com/achyar10/snmp-olt-zte/config"
	"github.com/achyar10/snmp-olt-zte/internal/utils"
	"github.com/redis/go-redis/v9"
)

var (
	redisHost               string
	redisPort               string
	redisPassword           string
	redisDB                 int
	redisMinIdleConnections int
	redisPoolSize           int
	redisPoolTimeout        int
)

func NewRedisClient(cfg *config.Config) *redis.Client {
	if os.Getenv("APP_ENV") == "development" || os.Getenv("APP_ENV") == "production" {
		redisHost = os.Getenv("REDIS_HOST")
		redisPort = os.Getenv("REDIS_PORT")
		redisPassword = os.Getenv("REDIS_PASSWORD")
		redisDB = utils.ConvertStringToInteger(os.Getenv("REDIS_DB"))
		redisMinIdleConnections = utils.ConvertStringToInteger(os.Getenv("REDIS_MIN_IDLE_CONNECTIONS"))
		redisPoolSize = utils.ConvertStringToInteger(os.Getenv("REDIS_POOL_SIZE"))
		redisPoolTimeout = utils.ConvertStringToInteger(os.Getenv("REDIS_POOL_TIMEOUT"))
	} else {
		redisHost = cfg.RedisCfg.Host
		redisPort = cfg.RedisCfg.Port
		redisPassword = cfg.RedisCfg.Password
		redisDB = cfg.RedisCfg.DB
		redisMinIdleConnections = cfg.RedisCfg.MinIdleConnections
		redisPoolSize = cfg.RedisCfg.PoolSize
		redisPoolTimeout = cfg.RedisCfg.PoolTimeout
	}

	return redis.NewClient(&redis.Options{
		Addr:         redisHost + ":" + redisPort,
		Password:     redisPassword,
		DB:           redisDB,
		MinIdleConns: redisMinIdleConnections,
		PoolSize:     redisPoolSize,
		PoolTimeout:  time.Duration(redisPoolTimeout) * time.Second,
	})
}
