package repository

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/k1v4/Pinger/backend/internal/entity"
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
	sql, args, err := cr.Builder.
		Select("*").
		Where(sq.Eq{"ip": ip}).
		ToSql()
	if err != nil {
		return entity.Container{}, fmt.Errorf("GetContainer: %w", err)
	}

	var container entity.Container

	err = cr.Pool.
		QueryRow(ctx, sql, args...).
		Scan(&container.IpAddr, &container.PingTime, &container.LastSuccessful)
	if err != nil {
		return entity.Container{}, fmt.Errorf("GetContainer: %w", err)
	}

	return container, nil
}

func (cr *ContainerRepo) GetAllContainers(ctx context.Context) ([]entity.Container, error) {
	sql, _, err := cr.Builder.
		Select("source, destination, original, translation").
		From("history").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("TranslationRepo - GetHistory - r.Builder: %w", err)
	}

	rows, err := cr.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("TranslationRepo - GetHistory - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	containers := make([]entity.Container, 0, _defaultEntityCap)

	for rows.Next() {
		container := entity.Container{}

		err = rows.Scan(&container.IpAddr, &container.PingTime, &container.LastSuccessful)
		if err != nil {
			return nil, fmt.Errorf("GetAllContainers: %w", err)
		}

		containers = append(containers, container)
	}

	return containers, nil
}

func (cr *ContainerRepo) AddContainer(ctx context.Context, ip string) (string, error) {
	sql, args, err := cr.Builder.
		Insert("containers").
		Columns("ip").
		Values(ip).
		ToSql()
	if err != nil {
		return "", fmt.Errorf("TranslationRepo - Store - r.Builder: %w", err)
	}

	_, err = cr.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return "", fmt.Errorf("TranslationRepo - Store - r.Pool.Exec: %w", err)
	}

	return ip, nil
}

func (cr *ContainerRepo) UpdateContainer(ctx context.Context, container entity.Container) (entity.Container, error) {
	sql, args, err := cr.Builder.Update("containers").
		Set("ping_time", container.PingTime).
		Set("last_successful", container.LastSuccessful).
		Where(sq.Eq{"ip": container.IpAddr}).
		ToSql()
	if err != nil {
		return entity.Container{}, err
	}

	_, err = cr.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return entity.Container{}, fmt.Errorf("UpdateContainer: %w", err)
	}

	return container, nil
}

func (cr *ContainerRepo) DeleteContainer(ctx context.Context, ip string) error {
	sql, args, err := cr.Builder.Delete("containers").Where(sq.Eq{"ip": ip}).ToSql()
	if err != nil {
		return fmt.Errorf("DeleteContainer: %w", err)
	}

	_, err = cr.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("DeleteContainer: %w", err)
	}

	return nil
}
