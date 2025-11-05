package middleware

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wfrscltech/vulcano/domain/mistake"
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

// @Description Define los detalles de una respuesta de error
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

// @Description Representa un error de cliente
type ClientError struct {
	// Resumen general del error
	Error string `json:"error"   example:"Invalid request"`
	// Menasje más específico del error
	Message string `json:"message" example:"El campo 'nombre' es obligatorio"`
}

// getType Obtiene el URI de la página de la especificación del código HTTP
func getType(status int) string {
	if t, ok := uriTypes[status]; ok {
		return baseDomain + t
	}

	return "about:blank"
}

// Helper para responder con Problem Details
func writeProblem(c echo.Context, code int, msg, instance string) error {
	c.Response().Header().Set("Content-Type", "application/problem+json")
	return c.JSON(code, ProblemDetails{
		Type:     getType(code),
		Title:    http.StatusText(code),
		Status:   code,
		Detail:   fmt.Sprint(msg),
		Instance: instance,
	})
}

// ProblemMiddleware Intercepta errores y los transforma a Problem Details
func ProblemMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err == nil {
			return nil
		}

		instance := fmt.Sprintf("[%s] %s", c.Request().Method, c.Request().RequestURI)

		var mk *mistake.Mistake
		if errors.As(err, &mk) {
			if mk.Code() == http.StatusInternalServerError {
				slog.Error("Handle Mistake", slog.String("error", err.Error()))
				return writeProblem(c, mk.Code(), mk.DevError(), instance)
			}

			return c.JSON(mk.Code(), ClientError{Error: http.StatusText(mk.Code()), Message: mk.Error()})
		} else {
			// Si el handler devolvió un *echo.HTTPError, lo convertimos
			if he, ok := err.(*echo.HTTPError); ok {
				slog.Error("Handle HTTP Error", slog.String("error", err.Error()))
				if he.Code >= http.StatusInternalServerError {
					return writeProblem(c, he.Code, he.Message.(string), instance)
				}

				return c.JSON(he.Code, ClientError{Error: http.StatusText(he.Code), Message: he.Message.(string)})
			} else {
				// Error inesperado o no capturado
				slog.Error("Handle Unknown Error", slog.String("error", err.Error()))
				return writeProblem(c, http.StatusInternalServerError, err.Error(), instance)
			}
		}
	}

}
