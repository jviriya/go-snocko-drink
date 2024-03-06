package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go-pentor-bank/internal/appresponse"
	"go-pentor-bank/internal/clog"
	"go-pentor-bank/internal/config"
	"net/http"
	"runtime"

	"github.com/rs/zerolog"
)

func Recover(c *fiber.Ctx) error {
	logger := clog.GetContextLog(c.Context())
	defer func(ctx *fiber.Ctx, log zerolog.Logger) {
		if rec := recover(); rec != nil {
			err, ok := rec.(error)
			if !ok {
				err = fmt.Errorf("%v", rec)
			}
			stack := make([]byte, 4<<10) // 4KB
			length := runtime.Stack(stack, false)

			log.Error().Stack().Err(err).Msgf("panic recover stack: %s", stack[:length])
			appresponse.JSONResponse(ctx, http.StatusInternalServerError, appresponse.IResponse{ErrorCode: config.EM.Internal.InternalServerError, Data: "load test panic"})
		}
	}(c, *logger)
	return c.Next()
}
