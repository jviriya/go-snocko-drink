package middleware

import (
	"github.com/gofiber/fiber/v2"
	"go-pentor-bank/internal/appresponse"
	"go-pentor-bank/internal/clog"
	"go-pentor-bank/internal/config"
	"go-pentor-bank/internal/pkg/gateway/cloudflare"
	"net/http"
	"strings"
)

func VerifyTurnstileCaptcha(c *fiber.Ctx) error {
	log := clog.GetContextLog(c.Context())

	URL := config.Conf.CloudFlare.Turnstile.URL
	secretKey := config.Conf.CloudFlare.Turnstile.SecretKey
	responseKey := strings.TrimSpace(c.Get("cf-turnstile-response"))
	remoteIP := strings.TrimSpace(c.Get("CF-Connecting-IP"))

	if responseKey == "" || remoteIP == "" {
		return appresponse.JSONResponse(c, http.StatusInternalServerError, appresponse.IResponse{ErrorCode: config.EM.Validation.ValidationFailed})
	}

	err := cloudflare.CallbackTurnstile(log, URL, secretKey, responseKey, remoteIP)
	if err != nil {
		return appresponse.JSONResponse(c, http.StatusInternalServerError, appresponse.IResponse{ErrorCode: config.EM.Internal.InternalServerError})
	}

	return c.Next()
}
