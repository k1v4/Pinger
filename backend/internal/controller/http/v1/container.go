package v1

import (
	"errors"
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

	// группа роутов для /v1/containers
	h := handler.Group("/containers")
	{
		// GET /v1/containers
		h.GET("/", r.AllContainers)

		// GET /v1/containers/{ip}
		h.GET("/:ip", r.Container)

		// POST /v1/containers
		h.POST("/:ip", r.CheckPingContainer)

		// PUT /v1/containers/{ip}
		h.PUT("/:ip", r.UpdateContainer)

		// DELETE /v1/containers/{ip}
		h.DELETE("/:ip", r.DeleteContainer)

	}
}

func (tr *conatainerRoutes) CheckPingContainer(c echo.Context) error {
	ip := c.Param("ip")
	ctx := c.Request().Context()

	u := new(dto.DtoPingContainer)
	if err := c.Bind(u); err != nil {
		tr.l.Error(ctx, fmt.Sprintf("http-v1-CheckPingContainer: %s", err))
		errorResponse(c, http.StatusBadRequest, "bad request")

		return fmt.Errorf("http-v1-CheckPingContainer: %w", err)
	}

	getContainer, err := tr.t.Container(ctx, ip)
	if err != nil {
		if errors.Is(err, usecase.ErrNoIp) {
			ip, err = tr.t.NewContainer(ctx, entity.Container{
				IpAddr:         ip,
				PingTime:       u.PingTime,
				LastSuccessful: u.LastSuccessful,
			})
			if err != nil {
				tr.l.Error(ctx, fmt.Sprintf("http-v1-CheckPingContainer: %s", err))
				errorResponse(c, http.StatusInternalServerError, "database problems")

				return fmt.Errorf("http-v1-CheckPingContainer: %w", err)
			}

			return c.JSON(http.StatusOK, dto.NewContainerResponse{Ip: ip})
		}

		tr.l.Error(ctx, fmt.Sprintf("http-v1-CheckPingContainer: %s", err))
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return fmt.Errorf("http-v1-CheckPingContainer: %w", err)
	}

	if !u.IsSuccessful {
		container, err2 := tr.t.UpdateContainer(ctx, entity.Container{
			IpAddr:         ip,
			PingTime:       u.PingTime,
			LastSuccessful: getContainer.LastSuccessful,
		})
		if err2 != nil {
			tr.l.Error(ctx, fmt.Sprintf("http-v1-CheckPingContainer: %s", err2))
			errorResponse(c, http.StatusInternalServerError, "database problems")

			return fmt.Errorf("http-v1-CheckPingContainer: %w", err2)
		}

		return c.JSON(http.StatusOK, container)
	}

	updContainer, err := tr.t.UpdateContainer(ctx, entity.Container{
		IpAddr:         ip,
		PingTime:       u.PingTime,
		LastSuccessful: u.LastSuccessful,
	})
	if err != nil {
		tr.l.Error(ctx, fmt.Sprintf("http-v1-CheckPingContainer: %s", err))
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return fmt.Errorf("http-v1-CheckPingContainer: %w", err)
	}

	return c.JSON(http.StatusOK, dto.NewContainerResponse{Ip: updContainer.IpAddr})
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

	u := new(dto.AddContainerRequest)
	if err := c.Bind(u); err != nil {
		tr.l.Error(ctx, fmt.Sprintf("http-v1-NewContainer: %s", err))
		errorResponse(c, http.StatusBadRequest, "bad request")

		return fmt.Errorf("http-v1-NewContainer: %w", err)
	}

	ip, err := tr.t.NewContainer(ctx, entity.Container{
		IpAddr:         ip,
		PingTime:       u.PingTime,
		LastSuccessful: u.LastSuccessful,
	})
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
