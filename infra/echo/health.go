package echo

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type healthHandler struct {
	version    string
	buildTime  string
	commitHash string
	start      int64
}

func NewHealthHandler(version, buildTime, commitHash string) *healthHandler {
	return &healthHandler{
		version:    version,
		buildTime:  buildTime,
		commitHash: commitHash,
		start:      time.Now().Unix(),
	}
}

// @Description	Define la respuesta de la API de healthcheck
type HealthResponse struct {
	// Estado del servidor
	Status string `json:"status"      example:"ok"`
	// Versión de la aplicación servidor
	Version string `json:"version"     example:"0.1.0"`
	// Hora de compilación de la aplicación
	BuildTime string `json:"build_time"  example:"2021-01-01 00:00:00"`
	// Hash del último commit de la aplicación
	CommitHash string `json:"commit_hash" example:"56abc123"`
	// Tiempo actual del servidor
	Time string `json:"time"        example:"2006-01-02 15:04:05"`
	// Tiempo desde el inicio del servidor
	Uptime string `json:"uptime"      example:"1h2m3s"`
}

// @Summary		Validación de funcionamiento
// @Description	Endpoint de validación de la API, indica si el servidor está en funcionamiento
// @Tags			Monitoring
// @Produce		json
// @Success		200	{object}	HealthResponse
// @Router			/health [get]
func (h *healthHandler) Healthcheck(c echo.Context) error {
	return c.JSON(
		http.StatusOK,
		HealthResponse{
			Status:     "ok",
			Version:    h.version,
			BuildTime:  h.buildTime,
			CommitHash: h.commitHash,
			Time:       time.Now().Format("2006-01-02 15:04:05"),
			Uptime:     time.Since(time.Unix(h.start, 0)).String(),
		},
	)
}
