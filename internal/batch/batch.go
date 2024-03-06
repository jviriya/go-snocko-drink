package batch

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/robfig/cron/v3"
	"go-pentor-bank/internal/appresponse"
	"go-pentor-bank/internal/clog"
	"go-pentor-bank/internal/config"
	"go-pentor-bank/internal/pkg/healthcheck"
	"go-pentor-bank/internal/routes"
	"net/http"
	"time"
)

var shutdown = false

func Batch(job []string, port string) {
	//log := clog.GetLog()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return appresponse.JSONResponse(ctx, http.StatusOK, appresponse.IResponse{
				ErrorCode: config.EM.Internal.InternalServerError,
				Error:     err,
			})
		},
	})

	//ctx, cancel := context.WithCancel(context.Background())
	//ctx := context.Background()
	cj := cron.New(cron.WithParser(cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)))
	for _, j := range job {
		switch j {
		case "test_batch":
			//go utils.SafeGoRoutines(func() { cronJob(10*time.Second, ctx,) })
		}
	}

	healthcheck.RunFiberHealthCheck(app)
	cj.Start()
	routes.RunningGracefully(app, port)
	cj.Stop()
}

func cronJob(duration time.Duration, ctx context.Context, job func(context.Context) bool) {
	log := clog.GetLog()
	for {

		if shutdown {
			return
		}

		start := time.Now()

		log.Info().Msg("Start")

		wait := job(ctx)

		log.Info().Msgf("Done %v", time.Since(start))

		if wait {

			nextRoundTime := time.Now().Add(duration)

			for nextRoundTime.After(time.Now()) {

				if shutdown {
					return
				}

				time.Sleep(time.Second)
			}
		}
	}
}
