package main

import (
	"context"
	"fmt"
	"github.com/k1v4/Pinger/backend/internal/config"
	"github.com/k1v4/Pinger/backend/internal/usecase"
	"github.com/k1v4/Pinger/backend/internal/usecase/repository"
	"github.com/k1v4/Pinger/backend/pkg/DB/postgres"
	"github.com/k1v4/Pinger/backend/pkg/logger"
	"github.com/labstack/echo/v4"
)

func main() {
	ctx := context.Background()
	loggerBack := logger.NewLogger()

	loggerBack.Info(ctx, "starting backend")

	cfg := config.MustLoadConfig()
	if cfg == nil {
		loggerBack.Error(ctx, "config is nil")
		return
	}

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBConfig.UserName,
		cfg.DBConfig.Password,
		cfg.DBConfig.Host,
		cfg.DBConfig.Port,
		cfg.DBConfig.DbName,
	)

	pg, err := postgres.New(url, postgres.MaxPoolSize(cfg.DBConfig.PoolMax))
	if err != nil {
		loggerBack.Error(ctx, fmt.Sprintf("app - Run - postgres.New: %s", err))
	}
	defer pg.Close()

	loggerBack.Info(ctx, "connected to database successfully")

	containerUseCase := usecase.New(
		repository.NewContainerRepo(pg),
	)
	_ = containerUseCase

	handler := echo.New()
	_ = handler
}
