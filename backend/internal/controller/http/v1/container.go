package v1

import (
	"fmt"
	"github.com/k1v4/Pinger/backend/internal/controller/dto"
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

	// Группа роутов для /v1/containers
	h := handler.Group("/containers")
	{
		// GET /v1/containers
		h.GET("/", r.AllContainers)

		// GET /v1/containers/{ip}
		h.GET("/:ip", r.Container)

		// POST /v1/containers
		h.POST("/:ip", r.NewContainer)

		// PUT /v1/containers/{ip}
		h.PUT("/:ip", r.UpdateContainer)

		// DELETE /v1/containers/{ip}
		h.DELETE("/:ip", r.DeleteContainer)

	}
}

func (tr *conatainerRoutes) Container(c echo.Context) error {
	ip := c.Param("ip")
	ctx := c.Request().Context()

	container, err := tr.t.Container(ctx, ip)
	if err != nil {
		tr.l.Error(ctx, fmt.Sprintf("http-v1-Container: %s", err))
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return fmt.Errorf("http-v1-Container: %w", err)
	}

	return c.JSON(http.StatusOK, container)
}

func (tr *conatainerRoutes) AllContainers(c echo.Context) error {
	ctx := c.Request().Context()

	containers, err := tr.t.AllContainers(ctx)
	if err != nil {
		tr.l.Error(ctx, fmt.Sprintf("http-v1-AllContainers: %s", err))
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return fmt.Errorf("http-v1-AllContainers: %w", err)
	}

	return c.JSON(http.StatusOK, containers)
}

func (tr *conatainerRoutes) NewContainer(c echo.Context) error {
	ip := c.Param("ip")
	ctx := c.Request().Context()

	ip, err := tr.t.NewContainer(ctx, ip)
	if err != nil {
		tr.l.Error(ctx, fmt.Sprintf("http-v1-NewContainer: %s", err))
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return fmt.Errorf("http-v1-NewContainer: %w", err)
	}

	return c.JSON(http.StatusOK, dto.NewContainerResponse{Ip: ip})
}

func (tr *conatainerRoutes) UpdateContainer(c echo.Context) error {
	ip := c.Param("ip")
	ctx := c.Request().Context()

	u := new(dto.UpdateContainerRequest)
	if err := c.Bind(u); err != nil {
		tr.l.Error(ctx, fmt.Sprintf("http-v1-UpdateContainer: %s", err))
		errorResponse(c, http.StatusBadRequest, "bad request")

		return fmt.Errorf("http-v1-UpdateContainer: %w", err)
	}

	container := entity.Container{
		IpAddr:         ip,
		PingTime:       u.PingTime,
		LastSuccessful: u.LastSuccessful,
	}

	container, err := tr.t.UpdateContainer(ctx, container)
	if err != nil {
		tr.l.Error(ctx, fmt.Sprintf("http-v1-UpdateContainer: %s", err))
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return fmt.Errorf("http-v1-UpdateContainer: %w", err)
	}

	return c.JSON(http.StatusOK, container)
}

func (tr *conatainerRoutes) DeleteContainer(c echo.Context) error {
	ctx := c.Request().Context()
	ip := c.Param("ip")

	err := tr.t.DeleteContainer(ctx, ip)
	if err != nil {
		tr.l.Error(ctx, fmt.Sprintf("http-v1-DeleteContainer: %s", err))
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return fmt.Errorf("http-v1-DeleteContainer: %w", err)
	}

	return c.JSON(http.StatusOK, dto.DeleteContainerResponse{IsSuccess: true})
}
