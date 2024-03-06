package middleware

import (
	"github.com/gofiber/fiber/v2"
	"go-pentor-bank/internal/appresponse"
	"go-pentor-bank/internal/common"
	"go-pentor-bank/internal/config"
	"net/http"
	"strings"
)

func VerifyPlatformOS(c *fiber.Ctx) error {
	platform := c.Get(common.PlatformContext)
	switch strings.TrimSpace(platform) {
	case common.AndroidPlatform:
		platform = common.AndroidPlatform
	case common.IosPlatform:
		platform = common.IosPlatform
	case common.WebPlatform:
		platform = common.WebPlatform
	default:
		return appresponse.JSONResponse(c, http.StatusUnauthorized, appresponse.IResponse{ErrorCode: config.EM.Auth.InvalidPlatform})
	}
	c.Locals(common.PlatformContext, platform)

	return c.Next()
}
