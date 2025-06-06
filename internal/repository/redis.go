package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/achyar10/snmp-olt-zte/internal/model"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// OnuRedisRepositoryInterface is an interface that represent the auth's repository contract
type OnuRedisRepositoryInterface interface {
	GetOnuIDCtx(ctx context.Context, key string) ([]model.OnuID, error)
	SetOnuIDCtx(ctx context.Context, key string, seconds int, onuId []model.OnuID) error
	DeleteOnuIDCtx(ctx context.Context, key string) error
	SaveONUInfoList(ctx context.Context, key string, seconds int, onuInfoList []model.ONUInfoPerBoard) error
	GetONUInfoList(ctx context.Context, key string) ([]model.ONUInfoPerBoard, error)
	GetOnlyOnuIDCtx(ctx context.Context, key string) ([]model.OnuOnlyID, error)
	SaveOnlyOnuIDCtx(ctx context.Context, key string, seconds int, onuId []model.OnuOnlyID) error
}

// Auth redis repository
type onuRedisRepo struct {
	redisClient *redis.Client
}

// NewOnuRedisRepo will create an object that represent the auth repository
func NewOnuRedisRepo(redisClient *redis.Client) OnuRedisRepositoryInterface {
	return &onuRedisRepo{redisClient}
}

// GetOnuIDCtx is a method to get onu id from redis
func (r *onuRedisRepo) GetOnuIDCtx(ctx context.Context, key string) ([]model.OnuID, error) {
	onuBytes, err := r.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get onu id from redis")
		return nil, errors.Wrap(err, "onuRedisRepo.GetOnuIDCtx.redisClient.Get")
	}

	var onuId []model.OnuID
	if err := json.Unmarshal(onuBytes, &onuId); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal onu id")
		return nil, errors.Wrap(err, "onuRedisRepo.GetOnuIDCtx.json.Unmarshal")
	}

	return onuId, nil
}

// SetOnuIDCtx is a method to set onu id to redis
func (r *onuRedisRepo) SetOnuIDCtx(ctx context.Context, key string, seconds int, onuId []model.OnuID) error {
	onuBytes, err := json.Marshal(onuId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal onu id")
		return errors.Wrap(err, "setRedisRepo.SetNewsCtx.json.Marshal")
	}

	if err := r.redisClient.Set(ctx, key, onuBytes, time.Second*time.Duration(seconds)).Err(); err != nil {
		log.Error().Err(err).Msg("Failed to set onu id to redis")
		return errors.Wrap(err, "onuRedisRepo.SetOnuIDCtx.redisClient.Set")
	}

	return nil
}

// DeleteOnuIDCtx is a method to delete onu id from redis
func (r *onuRedisRepo) DeleteOnuIDCtx(ctx context.Context, key string) error {
	if err := r.redisClient.Del(ctx, key).Err(); err != nil {
		log.Error().Err(err).Msg("Failed to delete onu id from redis")
		return errors.Wrap(err, "onuRedisRepo.DeleteOnuIDCtx.redisClient.Del")
	}

	return nil
}

// SaveONUInfoList is a method to save onu info list to redis
func (r *onuRedisRepo) SaveONUInfoList(
	ctx context.Context, key string, seconds int, onuInfoList []model.ONUInfoPerBoard,
) error {
	onuBytes, err := json.Marshal(onuInfoList)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal onu info list")
		return errors.Wrap(err, "onuRedisRepo.SaveONUInfoList.json.Marshal")
	}

	if err := r.redisClient.Set(ctx, key, onuBytes, time.Second*time.Duration(seconds)).Err(); err != nil {
		log.Error().Err(err).Msg("Failed to set onu info list to redis")
		return errors.Wrap(err, "onuRedisRepo.SaveONUInfoList.redisClient.Set")
	}

	return nil
}

// GetONUInfoList is a method to get onu info list from redis
func (r *onuRedisRepo) GetONUInfoList(ctx context.Context, key string) ([]model.ONUInfoPerBoard, error) {
	onuBytes, err := r.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get onu info list from redis")
		return nil, errors.Wrap(err, "onuRedisRepo.GetONUInfoList.redisClient.Get")
	}

	var onuInfoList []model.ONUInfoPerBoard
	if err := json.Unmarshal(onuBytes, &onuInfoList); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal onu info list")
		return nil, errors.Wrap(err, "onuRedisRepo.GetONUInfoList.json.Unmarshal")
	}

	return onuInfoList, nil
}

// GetOnlyOnuIDCtx is a method to get only onu id from redis
func (r *onuRedisRepo) GetOnlyOnuIDCtx(ctx context.Context, key string) ([]model.OnuOnlyID, error) {
	onuBytes, err := r.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get onu id from redis")
		return nil, errors.Wrap(err, "onuRedisRepo.GetOnlyOnuIDCtx.redisClient.Get")
	}

	var onuId []model.OnuOnlyID
	if err := json.Unmarshal(onuBytes, &onuId); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal onu id")
		return nil, errors.Wrap(err, "onuRedisRepo.GetOnlyOnuIDCtx.json.Unmarshal")
	}

	return onuId, nil
}

// SaveOnlyOnuIDCtx is a method to save only onu id to redis
func (r *onuRedisRepo) SaveOnlyOnuIDCtx(ctx context.Context, key string, seconds int, onuId []model.OnuOnlyID) error {
	onuBytes, err := json.Marshal(onuId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal onu id")
		return errors.Wrap(err, "onuRedisRepo.SaveOnlyOnuIDCtx.json.Marshal")
	}

	if err := r.redisClient.Set(ctx, key, onuBytes, time.Second*time.Duration(seconds)).Err(); err != nil {
		log.Error().Err(err).Msg("Failed to set onu id to redis")
		return errors.Wrap(err, "onuRedisRepo.SaveOnlyOnuIDCtx.redisClient.Set")
	}

	return nil
}
