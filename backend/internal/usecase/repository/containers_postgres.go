package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/k1v4/Pinger/backend/internal/entity"
	"github.com/k1v4/Pinger/backend/internal/usecase"
	"github.com/k1v4/Pinger/backend/pkg/DB/postgres"
)

const _defaultEntityCap = 64

type ContainerRepo struct {
	*postgres.Postgres
}

func NewContainerRepo(pg *postgres.Postgres) *ContainerRepo {
	return &ContainerRepo{
		Postgres: pg,
	}
}

func (cr *ContainerRepo) GetContainer(ctx context.Context, ip string) (entity.Container, error) {
	s, args, err := cr.Builder.
		Select("*").
		From("containers").
		Where(sq.Eq{"ip": ip}).
		ToSql()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Container{}, usecase.ErrNoIp
		}

		return entity.Container{}, fmt.Errorf("ContainerRepo-GetContainer: %w", err)
	}

	var container entity.Container

	err = cr.Pool.
		QueryRow(ctx, s, args...).
		Scan(&container.IpAddr, &container.PingTime, &container.LastSuccessful)
	if err != nil {
		//fmt.Println(err, sql.ErrNoRows)
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Container{}, usecase.ErrNoIp
		}

		return entity.Container{}, fmt.Errorf("ContainerRepo-GetContainer: %w", err)
	}

	return container, nil
}

func (cr *ContainerRepo) GetAllContainers(ctx context.Context) ([]entity.Container, error) {
	sql, _, err := cr.Builder.
		Select("*").
		From("containers").
		OrderBy("ip ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ContainerRepo-GetAllContainers-r.Builder: %w", err)
	}

	rows, err := cr.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("ContainerRepo-GetAllContainers-r.Pool.Query: %w", err)
	}
	defer rows.Close()

	containers := make([]entity.Container, 0, _defaultEntityCap)

	for rows.Next() {
		container := entity.Container{}

		err = rows.Scan(&container.IpAddr, &container.PingTime, &container.LastSuccessful)
		if err != nil {
			return nil, fmt.Errorf("ContainerRepo-GetAllContainers: %w", err)
		}

		containers = append(containers, container)
	}

	return containers, nil
}

func (cr *ContainerRepo) AddContainer(ctx context.Context, container entity.Container) (string, error) {
	sql, args, err := cr.Builder.
		Insert("containers").
		Columns("ip", "ping_time", "last_successful").
		Values(container.IpAddr, container.PingTime, container.LastSuccessful).
		ToSql()
	if err != nil {
		return "", fmt.Errorf("ContainerRepo-AddContainer: %w", err)
	}

	_, err = cr.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return "", fmt.Errorf("ContainerRepo-AddContainer: %w", err)
	}

	return container.IpAddr, nil
}

func (cr *ContainerRepo) UpdateContainer(ctx context.Context, container entity.Container) (entity.Container, error) {
	sql, args, err := cr.Builder.Update("containers").
		Set("ping_time", container.PingTime).
		Set("last_successful", container.LastSuccessful).
		Where(sq.Eq{"ip": container.IpAddr}).
		ToSql()
	if err != nil {
		return entity.Container{}, fmt.Errorf("ContainerRepo-UpdateContainer: %w", err)
	}

	_, err = cr.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return entity.Container{}, fmt.Errorf("ContainerRepo-UpdateContainer: %w", err)
	}

	return container, nil
}

func (cr *ContainerRepo) DeleteContainer(ctx context.Context, ip string) error {
	sql, args, err := cr.Builder.Delete("containers").Where(sq.Eq{"ip": ip}).ToSql()
	if err != nil {
		return fmt.Errorf("ContainerRepo-DeleteContainer: %w", err)
	}

	_, err = cr.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("ContainerRepo-DeleteContainer: %w", err)
	}

	return nil
}
