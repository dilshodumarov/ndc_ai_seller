package user

import (
	"sugurta/internal/entity"
	"context"
)

type User interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	Get(ctx context.Context, params map[string]string) (*entity.User, error)
	List(ctx context.Context, filter entity.UserFilter) ([]*entity.User, error)
	Update(ctx context.Context, investor *entity.User) error
	UpdatePassword(ctx context.Context, investor *entity.UpdatePasswordRequest) error
	Delete(ctx context.Context, guid string) error
	ListClients(ctx context.Context, filter entity.ClientFilter) (*entity.ListClients, error)
	GetClientByID(ctx context.Context, id string) (*entity.Client, error)
	GetByIDs(ctx context.Context, id string) ([]*entity.User, error)
	BlockUser(ctx context.Context, req entity.BlockUser) error
}
