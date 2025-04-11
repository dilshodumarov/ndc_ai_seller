package repository

import (
	"sugurta/internal/entity"
	"context"
)

type Business interface {
	Create(ctx context.Context, business *entity.CreateBusinessRequest) error
	Get(ctx context.Context, id string) (*entity.Business, error)
	List(ctx context.Context,  busness entity.GetAllBusinessesRequest) (*entity.GetAllBusinessesResponse, error)
	Update(ctx context.Context, role *entity.UpdateBusinessRequest) error
	Delete(ctx context.Context, id string) error
}
