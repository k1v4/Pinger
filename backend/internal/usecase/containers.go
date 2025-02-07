package usecase

import (
	"context"
	"fmt"
	"github.com/k1v4/Pinger/backend/internal/entity"
)

type ContainerUseCase struct {
	repo ContainerRepo
}

func New(r ContainerRepo) *ContainerUseCase {
	return &ContainerUseCase{
		repo: r,
	}
}

func (cus *ContainerUseCase) Container(ctx context.Context, ip string) (entity.Container, error) {
	container, err := cus.repo.GetContainer(ctx, ip)
	if err != nil {
		return entity.Container{}, fmt.Errorf("ContainerUseCase_Container: %w", err)
	}

	return container, nil
}

func (cus *ContainerUseCase) AllContainers(ctx context.Context) ([]entity.Container, error) {
	containers, err := cus.repo.GetAllContainers(ctx)
	if err != nil {
		return nil, fmt.Errorf("ContainerUseCase_AllContainers: %w", err)
	}

	return containers, nil
}

func (cus *ContainerUseCase) NewContainer(ctx context.Context, pingContainer entity.Container) (string, error) {
	ip, err := cus.repo.AddContainer(ctx, entity.Container{
		IpAddr:         pingContainer.IpAddr,
		PingTime:       pingContainer.PingTime,
		LastSuccessful: pingContainer.LastSuccessful,
	})
	if err != nil {
		return "", fmt.Errorf("ContainerUseCase_NewContainer: %w", err)
	}

	return ip, nil
}

func (cus *ContainerUseCase) UpdateContainer(ctx context.Context, container entity.Container) (entity.Container, error) {
	updateContainer, err := cus.repo.UpdateContainer(ctx, container)
	if err != nil {
		return entity.Container{}, fmt.Errorf("ContainerUseCase_UpdateContainer: %w", err)
	}

	return updateContainer, nil
}

func (cus *ContainerUseCase) DeleteContainer(ctx context.Context, ip string) error {
	err := cus.repo.DeleteContainer(ctx, ip)
	if err != nil {
		return fmt.Errorf("ContainerUseCase_DeleteContainer: %w", err)
	}

	return nil
}
