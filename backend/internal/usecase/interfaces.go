package usecase

import (
	"context"
	"github.com/k1v4/Pinger/backend/internal/entity"
)

type (
	Container interface {
		Container(ctx context.Context, ip string) (entity.Container, error)
		AllContainers(ctx context.Context) ([]entity.Container, error)
		NewContainer(ctx context.Context, pingContainer entity.Container) (string, error)
		UpdateContainer(ctx context.Context, container entity.Container) (entity.Container, error)
		DeleteContainer(ctx context.Context, ip string) error
		//Translate(context.Context, entity.Translation) (entity.Translation, error)
		//History(context.Context) ([]entity.Translation, error)
	}

	ContainerRepo interface {
		GetContainer(ctx context.Context, ip string) (entity.Container, error)
		GetAllContainers(ctx context.Context) ([]entity.Container, error)
		AddContainer(ctx context.Context, container entity.Container) (string, error)
		UpdateContainer(ctx context.Context, container entity.Container) (entity.Container, error)
		DeleteContainer(ctx context.Context, ip string) error
	}
)
