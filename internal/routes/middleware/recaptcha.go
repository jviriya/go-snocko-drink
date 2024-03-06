package middleware

import (
	"github.com/gofiber/fiber/v2"
	"go-pentor-bank/internal/appresponse"
	"go-pentor-bank/internal/clog"
	"go-pentor-bank/internal/config"
	"go-pentor-bank/internal/pkg/gateway/google"
	"net/http"
	"strings"
)

func VerifyReCaptcha(c *fiber.Ctx) error {
	log := clog.GetContextLog(c.Context())

	URL := config.Conf.Google.ReCaptchaV2.URL
	secretKey := config.Conf.Google.ReCaptchaV2.SecretKey
	responseKey := strings.TrimSpace(c.Get("g-recaptcha-response"))

	if responseKey == "" {
		return appresponse.JSONResponse(c, http.StatusInternalServerError, appresponse.IResponse{ErrorCode: config.EM.Validation.ValidationFailed})
	}

	err := google.CallbackVerifyReCaptcha(log, URL, secretKey, responseKey)
	if err != nil {
		return appresponse.JSONResponse(c, http.StatusInternalServerError, appresponse.IResponse{ErrorCode: config.EM.Internal.InternalServerError})
	}

	return c.Next()
}
