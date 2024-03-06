package healthcheck

import (
	"github.com/gofiber/fiber/v2"
	"go-pentor-bank/internal/appresponse"
	"go-pentor-bank/internal/config"
	"go-pentor-bank/internal/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func BatchHealthCheck() {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return appresponse.JSONResponse(ctx, http.StatusOK, appresponse.IResponse{
				ErrorCode: config.EM.Internal.InternalServerError,
				Error:     err,
			})
		},
	})

	RunFiberHealthCheck(app)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go utils.SafeGoRoutines(func() {
		_ = <-c
		_ = app.Shutdown()
	})

	if err := app.Listen(":5555"); err != nil {
		log.Fatalln(err)
	}
}
