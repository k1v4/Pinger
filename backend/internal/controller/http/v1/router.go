package v1

import (
	"github.com/k1v4/Pinger/backend/internal/usecase"
	"github.com/k1v4/Pinger/backend/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func NewRouter(handler *echo.Echo, l logger.Logger, t usecase.Container) {
	// Middleware
	handler.Use(middleware.Logger())
	handler.Use(middleware.Recover())

	handler.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	h := handler.Group("/v1")
	{
		newContainerRoutes(h, t, l)
	}
}
