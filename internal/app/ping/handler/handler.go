package handler

import (
	"github.com/gofiber/fiber/v2"
	"go-pentor-bank/internal/app/ping/service"
	"go-pentor-bank/internal/appresponse"
	"go-pentor-bank/internal/config"
	_ "image/png"
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

type Test struct {
	Name string `json:"name"`
}

func (h *PingHandler) Ping(c *fiber.Ctx) error {
	//h.svr.MigrateUserHasOpenAPI(c)
	msg := "[app]: " + config.Conf.ServerSetting.PingMessage
	return appresponse.JSONResponse(c, http.StatusOK, appresponse.IResponse{ErrorCode: config.EM.Success, Data: msg})
}

func (h *PingHandler) TestNoti(c *fiber.Ctx) error {
	//h.svr.MigrateOpenAPI(c)
	msg := "[app]: " + config.Conf.ServerSetting.PingMessage
	return appresponse.JSONResponse(c, http.StatusOK, appresponse.IResponse{ErrorCode: config.EM.Success, Data: msg})
}
