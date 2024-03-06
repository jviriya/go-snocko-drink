package routes

import (
	"github.com/gofiber/fiber/v2"
	"go-pentor-bank/internal/clog"
	"go-pentor-bank/internal/pkg/healthcheck"
	"go-pentor-bank/internal/utils"
	"os"
	"os/signal"
)

func RunningGracefully(app *fiber.App, port string) {
	log := clog.GetLog()
	go utils.SafeGoRoutines(func() {
		log.Print("Running at port :", port)
		if err := app.Listen(":" + port); err != nil {
			log.Panic().Err(err).Msg("Listen: ")
		}
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Kill, os.Interrupt)
	_ = <-c
	log.Info().Msg("Gracefully shutting down...")
	healthcheck.SetReadinessProbeStatusInternalServerError()
	err := app.Shutdown()
	if err != nil {
		log.Error().Err(err).Msg("app.Shutdown got err")
	}
}
