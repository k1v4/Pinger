package v1

import (
	"context"
	"github.com/k1v4/Pinger/backend/internal/entity"
	"github.com/k1v4/Pinger/backend/internal/usecase"
	"github.com/k1v4/Pinger/backend/pkg/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

type conatainerRoutes struct {
	t usecase.Container
	l logger.Logger
}

func newContainerRoutes(handler *echo.Group, t usecase.Container, l logger.Logger) {
	r := &conatainerRoutes{t, l}

	// Группа роутов для /translation
	h := handler.Group("/v1/containers")
	{
		// GET /translation/history
		h.GET("/", r.AllContainers)
	}
}

func (tr *conatainerRoutes) Container(ectx echo.Context) error {
	panic("implement me")
}

func (tr *conatainerRoutes) AllContainers(c echo.Context) error {
	c.JSON(http.StatusOK, "")
}

func (tr *conatainerRoutes) NewContainer(ectx echo.Context) {
	panic("implement me")
}

func (tr *conatainerRoutes) UpdateContainer(ectx echo.Context) {
	panic("implement me")
}

func (tr *conatainerRoutes) DeleteContainer(ectx echo.Context) {
	panic("implement me")
}
