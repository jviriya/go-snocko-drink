package api

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go-pentor-bank/internal/appresponse"
	"go-pentor-bank/internal/clog"
	"go-pentor-bank/internal/config"
	//"go-pentor-bank/internal/infra/mongodb"
	"go-pentor-bank/internal/pkg/healthcheck"
	//"go-pentor-bank/internal/repository/commondb"
	"go-pentor-bank/internal/routes"
	"go-pentor-bank/internal/routes/middleware"
	"strings"
)

func App(service, port, prefixPathFlag string) {
	log := clog.GetLog()

	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
		BodyLimit:   config.Conf.ServerSetting.MaxFileSize,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return appresponse.JSONResponse(ctx, fiber.StatusInternalServerError, appresponse.IResponse{
				ErrorCode: config.EM.Internal.InternalServerError,
				Error:     err,
			})
		},
	})

	allowOrigin := strings.Join(config.Conf.ServerSetting.AllowOrigins, ",")
	switch config.Conf.State {
	case config.StateDEV, config.StateSIT:
		allowOrigin = "*"
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigin,
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,HEAD,OPTIONS",
		AllowHeaders:     "Content-Type,Origin,Content-Length,Content-Label,cf-request-id,Authorization,Accept-Language,Platform,Captcha-ID,Captcha-Output,Gen-Time,Lot-Number,Pass-Token",
		AllowCredentials: true,
		MaxAge:           12 * 3600,
	}))

	healthcheck.RunFiberHealthCheck(app)

	routesV0 := NewRouterV0()
	var routesV1 []routes.Route

	routesV1 = NewRouterV1()
	prefixPath := "/api"
	if prefixPathFlag != "" {
		prefixPath = prefixPathFlag
	}

	// Set groupPath
	var groupPathV0, groupPathV1 string
	switch config.Conf.State {
	case config.StateDEV, config.StateSIT, config.StateUAT, config.StateProd, config.StateLoadTest:
		groupPathV1 = prefixPath + "/v1"
		groupPathV0 = prefixPath
	default:
		groupPathV1 = "/v1"
	}

	//rp := commondb.NewRepositoryV2(mongodb.MongoDBCon.SecureClient, mongodb.MongoDBCon.SecureClientShard)

	for _, ro := range routesV0 {
		if !ro.Test || config.Conf.State != config.StateProd {
			app.Add(ro.Method, groupPathV0+ro.Pattern, append(ro.Middleware, ro.Endpoint)...)
			log.Info().Msgf("[API] %v %v%v", ro.Method, groupPathV0, ro.Pattern)
		}
	}

	//v1
	v1 := app.Group(groupPathV1)
	v1.Use(middleware.AcceptLanguage())
	v1.Use(middleware.VerifyPlatformOS)
	v1.Use(middleware.LoggerWithRequestMeta())
	v1.Use(middleware.Recover)
	v1.Use(middleware.LogMetadataReqResp(service))
	//v1.Use(middleware.CheckVersion(rp))
	//v1.Use(middleware.CheckMaintenance(rp))

	for _, ro := range routesV1 {
		if !ro.Test || config.Conf.State != config.StateProd {
			v1.Add(ro.Method, ro.Pattern, append(ro.Middleware, ro.Endpoint)...)
			log.Info().Msgf("[API] %v %v%v", ro.Method, groupPathV1, ro.Pattern)
		}
	}

	routes.RunningGracefully(app, port)

	return
}

func NewRouterV0() []routes.Route {
	var routes []routes.Route
	routes = append(routes, ping()...)

	return routes
}

func NewRouterV1() []routes.Route {
	var routes []routes.Route
	//routes = append(routes, auth()...)

	return routes
}
