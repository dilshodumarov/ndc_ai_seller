package database

import (
	"context"
	"sugurta/internal/entity"
	"sugurta/internal/infrastructure/repository"
	"time"
)


type databaseService struct {
	ctxTimeout   time.Duration
	databaseRepo repository.Database
}

func NewDatabaseService(ctxTimeout time.Duration, repo repository.Database) Database {
	return &databaseService{
		ctxTimeout:   ctxTimeout,
		databaseRepo: repo,
	}
}

func (ds *databaseService) Create(ctx context.Context, req *entity.CreateDatabaseRequest) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, ds.ctxTimeout)
	defer cancel()

	return ds.databaseRepo.Create(ctx, req)
}

func (ds *databaseService) GetByID(ctx context.Context, guid string) (*entity.Database, error) {
	ctx, cancel := context.WithTimeout(ctx, ds.ctxTimeout)
	defer cancel()

	return ds.databaseRepo.GetByID(ctx, guid)
}

func (ds *databaseService) Update(ctx context.Context, req *entity.UpdateDatabaseRequest) error {
	ctx, cancel := context.WithTimeout(ctx, ds.ctxTimeout)
	defer cancel()

	return ds.databaseRepo.Update(ctx, req)
}

func (ds *databaseService) Delete(ctx context.Context, guid string) error {
	ctx, cancel := context.WithTimeout(ctx, ds.ctxTimeout)
	defer cancel()

	return ds.databaseRepo.Delete(ctx, guid)
}

func (ds *databaseService) List(ctx context.Context) ([]*entity.Database, error) {
	ctx, cancel := context.WithTimeout(ctx, ds.ctxTimeout)
	defer cancel()

	return ds.databaseRepo.List(ctx)
}
