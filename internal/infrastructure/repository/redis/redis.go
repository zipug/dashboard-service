package redis

import (
	"context"
	"dashboard/internal/common/service/config"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrPing          = errors.New("could not ping database")
	ErrKeyIsNotValid = errors.New("key is not valid")
)

type RedisRepository struct {
	rdb    *redis.Client
	config *redis.Options
	expire time.Duration
}

func NewRedisRepository(cfg *config.AppConfig) *RedisRepository {
	repo := &RedisRepository{}
	if err := repo.InvokeConnect(cfg); err != nil {
		e := fmt.Errorf("REDIS: redis://%s:%s@%s:%d/%d\n%w", cfg.Redis.User, cfg.Redis.Password, cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.DB, err)
		panic(e)
	}
	return repo
}

func (repo *RedisRepository) InvokeConnect(cfg *config.AppConfig) error {
	conf := redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.RedisPassword,
		DB:       cfg.Redis.DB,
		Protocol: 2,
	}
	rdb := redis.NewClient(&conf)
	repo.config = &conf
	repo.expire = cfg.OTP.ExpirationTime
	repo.rdb = rdb
	if err := repo.PingTest(); err != nil {
		panic(err)
	}
	return nil
}

func (repo *RedisRepository) PingTest() error {
	max_errs := 5
	errs := 0
	timeout := 1 * time.Second
	for max_errs > 0 {
		if err := repo.rdb.Ping(context.Background()).Err(); err != nil {
			fmt.Printf("could not ping database: %s\n", err.Error())
			fmt.Printf("retrying in %s\n", timeout)
			max_errs--
			errs++
			time.Sleep(timeout)
		}
		max_errs = 0
		errs = 0
	}
	if errs == 0 {
		return nil
	}
	return fmt.Errorf("%w: redis_uri: %s", ErrPing, repo.config.Addr)
}
