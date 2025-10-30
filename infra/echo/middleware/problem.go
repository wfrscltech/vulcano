package middleware

/*
Middleware para transformar errores en Problem Details segun el RFC 9457

Modo de uso:

	e := echo.New()
	// Usamos el middleware
	e.Use(ProblemMiddleware)

	// Ejemplo: recurso no encontrado
	e.GET("/item/:id", func(c echo.Context) error {
		id := c.Param("id")

		// Forzamos un 404 usando echo.HTTPError
		return echo.NewHTTPError(http.StatusNotFound, "Item "+id+" does not exist")
	})

	// Ejemplo: error interno simulado
	e.GET("/panic", func(c echo.Context) error {
		return fmt.Errorf("something went really wrong")
	})

	e.Logger.Fatal(e.Start(":8080"))

*/

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

const baseDomain = "https://developer.mozilla.org"

var uriTypes = map[int]string{
	http.StatusBadRequest:                    "/es/docs/Web/HTTP/Reference/Status/400",
	http.StatusUnauthorized:                  "/es/docs/Web/HTTP/Reference/Status/401",
	http.StatusPaymentRequired:               "/en-US/docs/Web/HTTP/Reference/Status/402",
	http.StatusForbidden:                     "/es/docs/Web/HTTP/Reference/Status/403",
	http.StatusNotFound:                      "/es/docs/Web/HTTP/Reference/Status/404",
	http.StatusMethodNotAllowed:              "/es/docs/Web/HTTP/Reference/Status/405",
	http.StatusNotAcceptable:                 "/en-US/docs/Web/HTTP/Reference/Status/406",
	http.StatusProxyAuthRequired:             "/en-US/docs/Web/HTTP/Reference/Status/407",
	http.StatusRequestTimeout:                "/es/docs/Web/HTTP/Reference/Status/408",
	http.StatusConflict:                      "/en-US/docs/Web/HTTP/Reference/Status/409",
	http.StatusGone:                          "/en-US/docs/Web/HTTP/Reference/Status/410",
	http.StatusLengthRequired:                "/en-US/docs/Web/HTTP/Reference/Status/411",
	http.StatusPreconditionFailed:            "/en-US/docs/Web/HTTP/Reference/Status/412",
	http.StatusRequestEntityTooLarge:         "/es/docs/Web/HTTP/Reference/Status/413",
	http.StatusRequestURITooLong:             "/en-US/docs/Web/HTTP/Reference/Status/414",
	http.StatusUnsupportedMediaType:          "/en-US/docs/Web/HTTP/Reference/Status/415",
	http.StatusRequestedRangeNotSatisfiable:  "/en-US/docs/Web/HTTP/Reference/Status/416",
	http.StatusExpectationFailed:             "/en-US/docs/Web/HTTP/Reference/Status/417",
	http.StatusTeapot:                        "/es/docs/Web/HTTP/Reference/Status/418",
	http.StatusMisdirectedRequest:            "/en-US/docs/Web/HTTP/Reference/Status/421",
	http.StatusUnprocessableEntity:           "/en-US/docs/Web/HTTP/Reference/Status/422",
	http.StatusLocked:                        "/en-US/docs/Web/HTTP/Reference/Status/423",
	http.StatusFailedDependency:              "/en-US/docs/Web/HTTP/Reference/Status/424",
	http.StatusTooEarly:                      "/en-US/docs/Web/HTTP/Reference/Status/425",
	http.StatusUpgradeRequired:               "/en-US/docs/Web/HTTP/Reference/Status/426",
	http.StatusPreconditionRequired:          "/en-US/docs/Web/HTTP/Reference/Status/428",
	http.StatusTooManyRequests:               "/en-US/docs/Web/HTTP/Reference/Status/429",
	http.StatusRequestHeaderFieldsTooLarge:   "/en-US/docs/Web/HTTP/Reference/Status/431",
	http.StatusUnavailableForLegalReasons:    "/en-US/docs/Web/HTTP/Reference/Status/451",
	http.StatusInternalServerError:           "/es/docs/Web/HTTP/Reference/Status/500",
	http.StatusNotImplemented:                "/en-US/docs/Web/HTTP/Reference/Status/501",
	http.StatusBadGateway:                    "/es/docs/Web/HTTP/Reference/Status/502",
	http.StatusServiceUnavailable:            "/es/docs/Web/HTTP/Reference/Status/503",
	http.StatusGatewayTimeout:                "/es/docs/Web/HTTP/Reference/Status/504",
	http.StatusHTTPVersionNotSupported:       "/es/docs/Web/HTTP/Reference/Status/505",
	http.StatusVariantAlsoNegotiates:         "/en-US/docs/Web/HTTP/Reference/Status/506",
	http.StatusInsufficientStorage:           "/en-US/docs/Web/HTTP/Reference/Status/507",
	http.StatusLoopDetected:                  "/en-US/docs/Web/HTTP/Reference/Status/508",
	http.StatusNotExtended:                   "/en-US/docs/Web/HTTP/Reference/Status/510",
	http.StatusNetworkAuthenticationRequired: "/en-US/docs/Web/HTTP/Reference/Status/511",
}

// ProblemDetails Define los detalles de una respuesta de error
type ProblemDetails struct {
	// URI que identifica el tipo de problema
	Type string `json:"type"`
	// Frase corta que resume el problema
	Title string `json:"title"`
	// El código HTTP de respuesta que acompaña al problema
	Status int `json:"status"`
	// Explicación más específica del problema
	Detail string `json:"detail"`
	// URI que identifica esta ocurrencia específica del problema
	Instance string `json:"instance"`
}

// getType Obtiene el URI de la página de la especificación del código HTTP
func getType(status int) string {
	if t, ok := uriTypes[status]; ok {
		return baseDomain + t
	}

	return "about:blank"
}

// Helper para responder con Problem Details
func writeProblem(c echo.Context, pd ProblemDetails) error {
	c.Response().Header().Set("Content-Type", "application/problem+json")
	return c.JSON(pd.Status, pd)
}

// ProblemMiddleware Intercepta errores y los transforma a Problem Details
func ProblemMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err == nil {
			return nil
		}

		instance := fmt.Sprintf("[%s] %s", c.Request().Method, c.Request().RequestURI)

		// Si el handler devolvió un *echo.HTTPError, lo convertimos
		if he, ok := err.(*echo.HTTPError); ok {
			pd := ProblemDetails{
				Type:     getType(he.Code),
				Title:    http.StatusText(he.Code),
				Status:   he.Code,
				Detail:   fmt.Sprint(he.Message),
				Instance: instance,
			}
			slog.Error("Handle HTTP Error", slog.String("error", err.Error()))
			return writeProblem(c, pd)
		}

		// Caso genérico: error inesperado
		pd := ProblemDetails{
			Type:     "about:blank",
			Title:    "Internal Server Error",
			Status:   http.StatusInternalServerError,
			Detail:   err.Error(),
			Instance: instance,
		}
		return writeProblem(c, pd)
	}
}
