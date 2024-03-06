package middleware

import (
	"github.com/gofiber/fiber/v2"
	"go-pentor-bank/internal/appresponse"
	"go-pentor-bank/internal/clog"
	"go-pentor-bank/internal/common"
	"go-pentor-bank/internal/config"
	"go-pentor-bank/internal/infra/redis"
	"go-pentor-bank/internal/utils"
	"net/http"
	"time"
)

type RequestParam struct {
	Phone string `json:"phone"`
}

func (m AuthMiddleware) VerifyRequestLimit(key string, limit int64, minute time.Duration, condition string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		//log := clog.GetContextLog(ctx.Context())
		log := clog.GetLog()
		var identification string
		switch condition {
		case common.IdIPRedisKey:
			identification = utils.GetIP(ctx)
		case common.PhoneRedisKey:
			var req RequestParam
			if err := ctx.BodyParser(&req); err != nil {
				log.Error().Err(err).Msg("parser error")
				return appresponse.JSONResponse(ctx, http.StatusOK, appresponse.IResponse{
					ErrorCode: config.EM.Internal.InternalServerError,
					Error:     err,
				})
			}
			identification = req.Phone
		}

		countKey := common.RateLimitRedisKey + key + "_" + identification
		expiredAtKey := common.RateLimitExpiredAtRedisKey + key + "_" + identification
		rs, err := redis.RedisClient.HIncrBy(ctx.Context(), countKey, "count", 1)
		if err != nil {
			log.Error().Err(err).Msg("rate limit HIncrBy got err")
			return appresponse.JSONResponse(ctx, http.StatusBadRequest, appresponse.IResponse{ErrorCode: config.EM.Internal.InternalServerError})
		}

		if rs == 1 { //first times
			expiredDuration := minute * time.Minute
			redis.RedisClient.Expire(ctx.Context(), countKey, expiredDuration)
			redis.RedisClient.Set(ctx.Context(), expiredAtKey, time.Now().Add(expiredDuration).Format(time.RFC3339Nano), expiredDuration)
		}

		if rs > limit {
			log.Error().Msgf("Too many request")
			expiredAtRaw, _ := redis.RedisClient.Get(ctx.Context(), expiredAtKey)
			expiredTime, _ := time.Parse(time.RFC3339Nano, expiredAtRaw)
			data := map[string]interface{}{
				"expiredTime": expiredTime,
			}
			return appresponse.JSONResponse(ctx, http.StatusBadRequest, appresponse.IResponse{ErrorCode: config.EM.Internal.RequestLimit, Data: data})
		}

		return ctx.Next()
	}
}
