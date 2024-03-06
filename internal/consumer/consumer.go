package consumer

import (
	"github.com/gofiber/fiber/v2"
	"go-pentor-bank/internal/appresponse"
	"go-pentor-bank/internal/clog"
	"go-pentor-bank/internal/config"
	"go-pentor-bank/internal/pkg/healthcheck"
	"go-pentor-bank/internal/utils"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Run(job, port, socketPortV2, socketPortV4 string) {
	log := clog.GetLog()

	//var shutdown bool

	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return appresponse.JSONResponse(ctx, http.StatusOK, appresponse.IResponse{
				ErrorCode: config.EM.Internal.InternalServerError,
				Error:     err,
			})
		},
	})
	switch job {
	case "test":
		//socketIO.RunBroadcast(socketPortV2, socketPortV4)
		//go utils.SafeGoRoutines(func() { test.NewConsumer(ctx, log, &shutdown) })
	}

	healthcheck.RunFiberHealthCheck(app)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	go utils.SafeGoRoutines(func() {
		_ = <-c
		_ = app.Shutdown()
	})

	if err := app.Listen(":" + port); err != nil {
		log.Fatal().Err(err).Send()
	}

	//shutdown = true
	utils.WaitGoRoutines()
}
