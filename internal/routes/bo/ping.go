package bo

import (
	"go-pentor-bank/internal/appbo/ping/handler"
	"go-pentor-bank/internal/appbo/ping/service"
	"go-pentor-bank/internal/routes"
	"net/http"
)

func ping() []routes.Route {
	//rp := repository.NewRepository(mongodb.MongoDBCon.DB, mongodb.MongoDBCon.SecureClient)
	pingHandler := newPingHandler()

	routes := []routes.Route{
		{
			Name:        "MakePing",
			Description: "MakePing",
			Method:      http.MethodGet,
			Pattern:     "/pingBo",
			Endpoint:    pingHandler.Ping,
		},
	}
	return routes
}

func newPingHandler() *handler.PingHandler {
	//jwtMgr := jwtManager.NewJWTManager(redis.RedisClient)

	pingService := service.NewPingService()
	pingHandler := handler.NewPingHandler(pingService)
	return pingHandler
}
