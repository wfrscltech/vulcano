package apidocs

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/swaggo/swag"
)

const htmlBody = `<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8" />
    <title>Documentación API REST EnterpriseBI</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
	<link rel="icon" href="/doc/static/favicon.png" type="image/x-icon" />
    <style>
      body {
        margin: 0;
        padding: 0;
        font-family: "Helvetica", sans-serif;
      }
      header {
        background: #fafafa;
        padding: 10px 20px;
        display: flex;
        align-items: center;
      }
      header img {
        height: 40px;
      }
      redoc {
        display: block;
        height: calc(100vh - 60px);
      }
    </style>
  </head>
  <body>
    <header>
      <img src="/doc/static/logo.png" alt="Logo CLTech" /> <strong>Documentación API REST EnterpriseBI</strong>
    </header>
    <redoc spec-url="/doc/spec/swagger.json"></redoc>
    <script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"></script>
  </body>
</html>`

//go:embed static/*
var staticFiles embed.FS

// APIDocsManager Inicializa el endpoint de la documentación de la API
func APIDocsManager(e *echo.Echo, swinfo *swag.Spec) {
	fsys, err := fs.Sub(staticFiles, "static")
	if err != nil {
		slog.Error("Error al cargar recursos estáticos", slog.String("error", err.Error()))
	}

	fileServer := http.FileServer(http.FS(fsys))
	e.GET("/doc/static/*", echo.WrapHandler(http.StripPrefix("/doc/static/", fileServer)))

	e.GET("/doc/api", func(c echo.Context) error {
		return c.HTML(http.StatusOK, htmlBody)
	})

	e.GET("/doc/spec/swagger.json", func(c echo.Context) error {
		return c.JSONBlob(http.StatusOK, []byte(swinfo.ReadDoc()))
	})
}
