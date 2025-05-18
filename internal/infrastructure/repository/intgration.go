package repository

import (
	"context"
	"sugurta/internal/entity"
)

type Integration interface {
	Create(ctx context.Context, req *entity.IntegrationCreate) error
	Update(ctx context.Context, req *entity.IntegrationUpdate) (*entity.IntegrationUpdateResponse, error)
	Delete(ctx context.Context, id string) error
	GetByOwnerID(ctx context.Context, req *entity.IntegrationRequest) (*entity.IntegrationGetResponse, error)
	UpdateStatus(ctx context.Context, req *entity.IntegrationUpdateStatus) (*entity.IntegrationUpdateStatusResponse, error)
	GetTokenUsageList(ctx context.Context, req *entity.IntegrationListRequest) (*entity.IntegrationListResponse, error)
	CheckIntegrationExistence(ctx context.Context, businessID string) (*entity.IntegrationExistenceResponse, error)
}
