package database

import (
	"context"
	"sugurta/internal/entity"
)

type Database interface {
	Create(ctx context.Context, req *entity.CreateDatabaseRequest) (string, error)
	GetByID(ctx context.Context, guid string) (*entity.Database, error)
	Update(ctx context.Context, req *entity.UpdateDatabaseRequest) error
	Delete(ctx context.Context, guid string) error
	List(ctx context.Context) ([]*entity.Database, error)
}
