package clienttype

import (
	"sugurta/internal/entity"
	"context"
)

type ClientType interface {
	Create(ctx context.Context, data *entity.CreateClientTypeRequest) error
	Get(ctx context.Context, id string) (*entity.ClientType, error)
	List(ctx context.Context, limit, page int) ([]*entity.ClientType, error)
	Update(ctx context.Context, role *entity.UpdateClientType) error
	Delete(ctx context.Context, id string) error
}
