package redis

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/rs/zerolog"
	"go-pentor-bank/internal/config"
	"strings"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

const (
	PrefixCache         = "cache"
	CacheDuration       = time.Duration(5 * time.Minute)
	CacheDurationV2     = time.Duration(1 * time.Second)
	CachePlayerDulation = time.Duration(24 * time.Hour)
)

var RedisClient *Redis

type Redis struct {
	Client *goredis.Client
}

func NewRedisClient(log *zerolog.Logger) error {
	log.Info().Msg("Connecting to Redis..")

	client, err := dialToRedis(config.Conf.Redis.Common.Addr, config.Conf.Redis.Common.Password, config.Conf.Redis.Common.Port)
	if err != nil {
		return err
	}
	RedisClient = &Redis{
		Client: client,
	}

	log.Info().Msg("Connecting to Redis success!!")
	return nil
}

func GetRedisPrefix() string {
	return config.Conf.Redis.RedisPrefix
}

func (r *Redis) Set(c context.Context, keyRaw string, data interface{}, duration time.Duration) error {
	var dataByte []byte
	if dstr, valid := data.(string); valid {
		dataByte = []byte(dstr)
	} else {
		var err error
		dataByte, err = sonic.Marshal(data)
		if err != nil {
			return err
		}
	}

	key := fmt.Sprintf("%s_%s", config.Conf.Redis.RedisPrefix, keyRaw)
	if err := r.Client.Set(c, key, string(dataByte), duration).Err(); err != nil {
		return err
	}
	return nil
}

func (r *Redis) SetV2(c context.Context, keyRaw string, data interface{}, duration time.Duration) (string, error) {
	var dataByte []byte
	if dstr, valid := data.(string); valid {
		dataByte = []byte(dstr)
	} else {
		var err error
		dataByte, err = sonic.Marshal(data)
		if err != nil {
			return "", err
		}
	}

	key := fmt.Sprintf("%s_%s", config.Conf.Redis.RedisPrefix, keyRaw)
	rs, err := r.Client.Set(c, key, string(dataByte), duration).Result()
	if err != nil {
		return "", err
	}
	return rs, nil
}

func (r *Redis) HSet(c context.Context, keyRaw string, values ...interface{}) (int64, error) {
	key := fmt.Sprintf("%s_%s", config.Conf.Redis.RedisPrefix, keyRaw)

	var data []interface{}
	for i := 0; i < len(values); i = i + 2 {
		val := values[i+1]
		var dataByte []byte
		if dstr, valid := val.(string); valid {
			dataByte = []byte(dstr)
		} else {
			var err error
			dataByte, err = sonic.Marshal(val)
			if err != nil {
				return 0, err
			}
		}
		data = append(data, values[i], string(dataByte))
	}
	result, err := r.Client.HSet(c, key, data).Result()
	if err != nil {
		return result, err
	}
	return result, nil

}

func (r *Redis) Parse(c context.Context, keyRaw string, data interface{}) error {
	key := fmt.Sprintf("%s_%s", config.Conf.Redis.RedisPrefix, keyRaw)
	dataStr, err := r.Client.Get(c, key).Result()
	if err != nil && err != goredis.Nil {
		return err
	}
	if dataStr == "" {
		return goredis.Nil
	}

	err = sonic.Unmarshal([]byte(dataStr), data)
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) GetPrefixKey() string {
	return fmt.Sprintf("%s_", config.Conf.Redis.RedisPrefix)
}

func (r *Redis) Exists(c context.Context, keysRaw ...string) (int64, error) {
	keys := []string{}
	for _, v := range keysRaw {
		keys = append(keys, fmt.Sprintf("%s_%s", config.Conf.Redis.RedisPrefix, v))
	}

	data, err := r.Client.Exists(c, keys...).Result()
	if err != nil && err != goredis.Nil {
		return 0, err
	}

	return data, nil
}

func (r *Redis) Get(c context.Context, keyRaw string) (string, error) {
	key := fmt.Sprintf("%s_%s", config.Conf.Redis.RedisPrefix, keyRaw)
	data, err := r.Client.Get(c, key).Result()
	if err != nil && err != goredis.Nil {
		return "", err
	}
	if data == "" {
		return "", goredis.Nil
	}
	return data, nil
}

func (r *Redis) HGet(c context.Context, keyRaw, field string) (string, error) {
	key := fmt.Sprintf("%s_%s", config.Conf.Redis.RedisPrefix, keyRaw)
	data, err := r.Client.HGet(c, key, field).Result()
	if err != nil {
		return "", err
	}
	return data, nil
}

func (r *Redis) HParse(c context.Context, keyRaw, field string, data interface{}) error {
	key := fmt.Sprintf("%s_%s", config.Conf.Redis.RedisPrefix, keyRaw)
	dataStr, err := r.Client.HGet(c, key, field).Result()
	if err != nil && err != goredis.Nil {
		return err
	}
	if dataStr == "" {
		return goredis.Nil
	}

	err = sonic.Unmarshal([]byte(dataStr), data)
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) HMGet(c context.Context, keyRaw string, fields ...string) ([]string, error) {
	key := fmt.Sprintf("%s_%s", config.Conf.Redis.RedisPrefix, keyRaw)
	data, err := r.Client.HMGet(c, key, fields...).Result()
	if err != nil && err != goredis.Nil {
		return nil, err
	}

	rs := []string{}
	for _, v := range data {
		rs = append(rs, v.(string))
	}

	return rs, nil
}

//func (r *Redis) HMParse(c context.Context, keyRaw string, field []string, data ...*interface{}) error {
//	fmt.Println(data[0])
//	fmt.Println(data[1])
//
//	if len(field) != len(data) {
//		return errors.New("invalid parameters")
//	}
//
//	key := fmt.Sprintf("%s_%s", config.Conf.Redis.RedisPrefix, keyRaw)
//	aryDataStr, err := r.Client.HMGet(c, key, field...).Result()
//	if err != nil && err != goredis.Nil {
//		return err
//	}
//
//	for i, row := range aryDataStr {
//		val := data[i]
//		if row == nil {
//			continue
//		} else {
//			err = sonic.Unmarshal([]byte(row.(string)), val)
//			if err != nil {
//				return err
//			}
//		}
//	}
//
//	return nil
//}

func (r *Redis) HGetAll(c context.Context, keyRaw string) (map[string]string, error) {
	key := fmt.Sprintf("%s_%s", config.Conf.Redis.RedisPrefix, keyRaw)
	data, err := r.Client.HGetAll(c, key).Result()
	if err != nil {
		return map[string]string{}, err
	}
	return data, nil
}

func (r *Redis) HDel(c context.Context, keyRaw string, fields ...string) (int64, error) {
	key := fmt.Sprintf("%s_%s", config.Conf.Redis.RedisPrefix, keyRaw)
	rs, err := r.Client.HDel(c, key, fields...).Result()
	if err != nil {
		return 0, err
	}
	return rs, nil
}

func (r *Redis) Del(c context.Context, keyRaw ...string) (int64, error) {
	delKey := []string{}
	for _, v := range keyRaw {
		delKey = append(delKey, fmt.Sprintf("%v_%v", config.Conf.Redis.RedisPrefix, v))
	}
	rs, err := r.Client.Del(c, delKey...).Result()
	if err != nil {
		return 0, err
	}
	return rs, nil
}

func (r *Redis) DelRaw(c context.Context, addPrefix bool, keyRaw ...string) error {
	delKey := []string{}
	for _, v := range keyRaw {
		k := v
		if addPrefix {
			k = fmt.Sprintf("%s_%s", config.Conf.Redis.RedisPrefix, k)
		}
		delKey = append(delKey, k)
	}

	_, err := r.Client.Del(c, delKey...).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) Expire(c context.Context, keyRaw string, duration time.Duration) error {
	key := fmt.Sprintf("%s_%s", config.Conf.Redis.RedisPrefix, keyRaw)
	err := r.Client.Expire(c, key, duration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) MGet(c context.Context, keysRaw []string) ([]interface{}, error) {
	keys := []string{}
	for _, v := range keysRaw {
		keys = append(keys, fmt.Sprintf("%s_%s", config.Conf.Redis.RedisPrefix, v))
	}
	data, err := r.Client.MGet(c, keys...).Result()
	if err != nil && err != goredis.Nil {
		return []interface{}{}, err
	}
	return data, nil
}

func (r *Redis) Keys(c context.Context, keyRaw string) ([]string, error) {
	key := fmt.Sprintf("%s_%s", config.Conf.Redis.RedisPrefix, keyRaw)
	listRaw, err := r.Client.Keys(c, key).Result()
	if err != nil {
		return nil, err
	}

	list := []string{}
	for _, v := range listRaw {
		list = append(list, strings.TrimPrefix(v, config.Conf.Redis.RedisPrefix+"_"))
	}
	return list, nil
}

func (r *Redis) IncrBy(c context.Context, keyRaw string, value int64) (int64, error) {
	key := fmt.Sprintf("%s_%s", config.Conf.Redis.RedisPrefix, keyRaw)
	result, err := r.Client.IncrBy(c, key, value).Result()
	if err != nil {
		return result, err
	}
	return result, nil
}

func (r *Redis) IncrByFloat(c context.Context, keyRaw string, value float64) (float64, error) {
	key := fmt.Sprintf("%s_%s", config.Conf.Redis.RedisPrefix, keyRaw)
	result, err := r.Client.IncrByFloat(c, key, value).Result()
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (r *Redis) HIncrByFloat(c context.Context, keyRaw, field string, value float64) (float64, error) {
	key := fmt.Sprintf("%s_%s", config.Conf.Redis.RedisPrefix, keyRaw)
	rs, err := r.Client.HIncrByFloat(c, key, field, value).Result()
	if err != nil {
		return rs, err
	}
	return rs, nil
}

func (r *Redis) HIncrBy(c context.Context, keyRaw, field string, value int64) (int64, error) {
	key := fmt.Sprintf("%s_%s", config.Conf.Redis.RedisPrefix, keyRaw)
	rs, err := r.Client.HIncrBy(c, key, field, value).Result()
	if err != nil {
		return rs, err
	}
	return rs, nil
}

func (r *Redis) SetFloat(c context.Context, keyRaw string, value float64, duration time.Duration) error {
	key := fmt.Sprintf("%s_%s", config.Conf.Redis.RedisPrefix, keyRaw)
	_, err := r.Client.Set(c, key, value, duration).Result()
	if err != nil {
		return err
	}
	return nil
}

func dialToRedis(addr, password string, port int) (*goredis.Client, error) {
	client := goredis.NewClient(&goredis.Options{
		Addr:            fmt.Sprintf("%s:%d", addr, port), //redis port
		Password:        password,
		DialTimeout:     time.Duration(config.Conf.Redis.DialTimeout) * time.Second,
		ReadTimeout:     time.Duration(config.Conf.Redis.ReadTimeout) * time.Second,
		WriteTimeout:    time.Duration(config.Conf.Redis.WriteTimeout) * time.Second,
		PoolSize:        config.Conf.Redis.PoolSize,
		PoolTimeout:     time.Duration(config.Conf.Redis.PoolTimeout) * time.Second,
		MaxIdleConns:    config.Conf.Redis.MaxIdleConns,
		MaxActiveConns:  0,
		ConnMaxLifetime: time.Duration(config.Conf.Redis.MaxConnLifetime) * time.Second,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
