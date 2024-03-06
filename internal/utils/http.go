package utils

import (
	"github.com/bytedance/sonic"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GenerateHeaderJson() map[string]string {
	//requestid := GetRequestID(c)
	//appid := GetAppID(c)
	header := map[string]string{
		"Content-Type": "application/json",
		//"request-id":     requestid,
		//"request-app-id": appid,
	}
	return header
}

func GenerateHeaderEmtpy(c *gin.Context) map[string]string {
	header := map[string]string{}
	return header
}

func GenerateHeaderJsonOri(c *gin.Context) map[string]string {
	header := map[string]string{
		"Content-Type": "application/json",
	}
	return header
}

func GenerateHeaderUrlencoded() map[string]string {
	header := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	return header
}

func GenerateHeaderPlaintext(c *gin.Context) map[string]string {
	header := map[string]string{
		"Content-Type": "text/plain",
	}
	return header
}

func GenerateHeaderUrlencoded2(c *gin.Context) map[string]string {
	header := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	return header
}

func GetRequestID(c *gin.Context) string {
	requestid, _ := c.Request.Context().Value("request-id").(string)
	return requestid
}

func GetAppID(c *gin.Context) string {
	appid, _ := c.Request.Context().Value("request-app-id").(string)
	return appid
}

func BalanceTimeOut(underMaintenance []bool) *http.Client {
	var timeOut time.Duration = 10
	if len(underMaintenance) != 0 && underMaintenance[0] {
		timeOut = 5
	}

	return &http.Client{
		Timeout: timeOut * time.Second,
	}
}

func SetClientTimeOut(duration time.Duration) *http.Client {
	return &http.Client{
		Timeout: duration * time.Second,
	}
}

func ToJsonString(data interface{}) string {
	resp, _ := sonic.MarshalString(data)
	return resp
}
