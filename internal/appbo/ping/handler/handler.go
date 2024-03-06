package handler

import (
	"github.com/gofiber/fiber/v2"
	"go-pentor-bank/internal/appbo/ping/service"
	"go-pentor-bank/internal/appresponse"
	"go-pentor-bank/internal/config"
	"net/http"
)

type PingHandler struct {
	svr *service.PingService
}

func NewPingHandler(pingService *service.PingService) *PingHandler {
	return &PingHandler{
		svr: pingService,
	}
}

func (h *PingHandler) Ping(c *fiber.Ctx) error {
	msg := "[bo]: " + config.Conf.ServerSetting.PingMessage
	return appresponse.JSONResponse(c, http.StatusOK, appresponse.IResponse{ErrorCode: config.EM.Success, Data: msg})
}
