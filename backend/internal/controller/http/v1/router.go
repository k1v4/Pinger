package v1

import (
	"github.com/k1v4/Pinger/backend/internal/usecase"
	"github.com/k1v4/Pinger/backend/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(handler *echo.Echo, l logger.Logger, t usecase.Container) {
	// Middleware
	handler.Use(middleware.Logger())
	handler.Use(middleware.Recover())

	// Swagger
	// handler.GET("/swagger/*", echoSwagger.WrapHandler(swaggerFiles.Handler))

	// K8s probe
	//handler.GET("/healthz", func(c echo.Context) error {
	//	return c.NoContent(http.StatusOK)
	//})

	// Prometheus metrics
	// handler.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	// Routers
	h := handler.Group("/v1")
	{
		newContainerRoutes(h, t, l)
	}
}
