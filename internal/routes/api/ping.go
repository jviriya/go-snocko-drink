package api

import (
	"github.com/gofiber/fiber/v2"
	"go-pentor-bank/internal/app/ping/handler"
	"go-pentor-bank/internal/app/ping/service"
	"go-pentor-bank/internal/infra/mongodb"
	"go-pentor-bank/internal/infra/redis"
	"go-pentor-bank/internal/pkg/jwtManager"
	"go-pentor-bank/internal/repository"
	"go-pentor-bank/internal/routes"
	"net/http"
)

func ping() []routes.Route {
	rp := repository.NewCommonDBRepository(mongodb.MongoDBCon.SecureClient, mongodb.MongoDBCon.SecureClientShard)
	pingHandler := newPingHandler(rp)

	routes := []routes.Route{
		{
			Name:        "MakePing",
			Description: "MakePing",
			Method:      http.MethodGet,
			Pattern:     "/ping",
			Endpoint:    pingHandler.Ping,
			Middleware:  []fiber.Handler{},
		},
	}
	return routes
}

func newPingHandler(rp *repository.CommonDBRepository) *handler.PingHandler {
	pingService := service.NewPingService(rp.CommonDBRepo, jwtManager.NewJWTManager(redis.RedisClient))
	pingHandler := handler.NewPingHandler(pingService)
	return pingHandler
}
