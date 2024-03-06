package doc

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go-pentor-bank/internal/appresponse"
	"go-pentor-bank/internal/config"
	"go-pentor-bank/internal/pkg/healthcheck"
	"go-pentor-bank/internal/routes"
	"go-pentor-bank/internal/routes/middleware"
	"strings"
)

func App(service, port string) {

	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return appresponse.JSONResponse(ctx, fiber.StatusInternalServerError, appresponse.IResponse{
				ErrorCode: config.EM.Internal.InternalServerError,
				Error:     err,
			})
		},
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(config.Conf.ServerSetting.AllowOrigins, ","),
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,HEAD,OPTIONS",
		AllowHeaders:     "Origin,Content-Length,Content-Label,cf-request-id,Authorization,Accept-Language,Platform",
		AllowCredentials: true,
		MaxAge:           12 * 3600,
	}))

	healthcheck.RunFiberHealthCheck(app)

	prefixPath := "/doc"

	// Set groupPath
	var groupPath string
	switch config.Conf.State {
	case config.StateDEV, config.StateSIT, config.StateUAT, config.StateProd:
		groupPath = prefixPath + "/v1"
	default:
		groupPath = "/v1"
	}

	v1 := app.Group(groupPath)

	v1.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			"pbank": "FoxGopherDev",
		},
	}))

	v1.Use(middleware.AcceptLanguage())
	v1.Use(middleware.LoggerWithRequestMeta())
	v1.Use(middleware.Recover)
	v1.Use(middleware.LogMetadataReqResp(service))

	v1.Static("", "./docs")

	routes.RunningGracefully(app, port)

	return
}
