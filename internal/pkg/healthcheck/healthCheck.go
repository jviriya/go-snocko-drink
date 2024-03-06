package healthcheck

import (
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	"github.com/gorilla/mux"
	"net/http"
)

func RunFiberHealthCheck(app *fiber.App) {
	app.Add(http.MethodGet, "/healthz", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON("healthy")
	})
	app.Add(http.MethodGet, "/status", func(c *fiber.Ctx) error {
		return c.Status(readinessProbeStatus).JSON("healthy")
	})
}

func RunMuxHealthCheck(r *mux.Router) {
	r.HandleFunc("/healthz", func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("healthy"))
	}).Methods(http.MethodGet)

	r.HandleFunc("/status", func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(readinessProbeStatus)
		w.Write([]byte("healthy"))
	}).Methods(http.MethodGet)
}

func RunGinHealthCheck(r *gin.Engine) {
	r.GET("/healthz", func(context *gin.Context) {
		context.JSON(http.StatusOK, "healthy")
	})

	r.GET("/status", func(context *gin.Context) {
		context.JSON(readinessProbeStatus, "healthy")
	})
}

func RunHttpMuxHealthCheck(r *http.ServeMux) {
	r.HandleFunc("/healthz", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})

	r.HandleFunc("/status", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(readinessProbeStatus)
	})
}

func RunHttpHealthCheck(r *http.ServeMux) {
	r.HandleFunc("/healthz", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})

	r.HandleFunc("/status", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(readinessProbeStatus)
	})
}

func HttpStatus(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
}

func HttpHealthz(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(readinessProbeStatus)
}
