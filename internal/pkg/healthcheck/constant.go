package healthcheck

import "net/http"

var (
	readinessProbeStatus = http.StatusOK
)

func SetReadinessProbeStatusInternalServerError() {
	readinessProbeStatus = http.StatusInternalServerError
}
