package commondb

import (
	"context"
	"go-pentor-bank/internal/clog"
	"go-pentor-bank/internal/common"
	"go-pentor-bank/internal/config"
	"go-pentor-bank/internal/infra/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	IosPlatformUiText     UiTextPlatform = common.IosPlatform
	AndroidPlatformUiText UiTextPlatform = common.AndroidPlatform
)

type (
	UiTextPlatform string
	UiText         struct {
		ID       primitive.ObjectID `bson:"_id" json:"id"`
		Code     string             `bson:"code" json:"code"`
		Platform string             `bson:"platform" json:"platform"`
		Lang     config.Lang        `bson:"lang" json:"lang"`
	}
)

func (rp *SecureRepository) InsertUiText(ctx context.Context) {

}

func (rp *SecureRepository) DeleteManyUiText(ctx context.Context, platform UiTextPlatform) (int64, error) {
	filter := bson.M{
		"platform": platform,
	}
	rs, err := rp.Secure.DB.Collection(cltUiText).DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	return rs.DeletedCount, nil
}

func (rp *SecureRepository) BulkUiText(ctx context.Context, writes []mongo.WriteModel) (int64, error) {
	resp, err := rp.SecureShard.DB.Collection(cltUiText).BulkWrite(ctx, writes)
	if err != nil {
		return 0, err
	}
	return resp.ModifiedCount, nil
}

func (rp *SecureRepository) FindOneUiTextByPlatform(ctx context.Context, platform UiTextPlatform) ([]UiText, error) {
	var data []UiText

	key := common.CacheUiTextRedisKey + string(platform)
	err := redis.RedisClient.Parse(ctx, key, &data)
	if err == nil {
		return data, nil
	}

	filter := bson.M{
		"platform": platform,
	}

	cur, err := rp.Secure.DB.Collection(cltUiText).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	err = cur.All(ctx, &data)
	if err != nil {
		return nil, err
	}

	err = redis.RedisClient.Set(ctx, key, data, 1*time.Hour)
	if err != nil {
		log := clog.GetLog()
		log.Error().Err(err).Msg("RedisClient.Set got err")
	}

	return data, nil
}
