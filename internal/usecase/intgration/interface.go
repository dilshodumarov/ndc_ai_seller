package integration

import (
	"context"
	"sugurta/internal/entity"
)

type Integration interface {
	Create(ctx context.Context, req *entity.IntegrationCreate) error
	Update(ctx context.Context, req *entity.IntegrationUpdate) (*entity.IntegrationUpdateResponse, error)
	UpdateStatus(ctx context.Context, req *entity.IntegrationUpdateStatus) (*entity.IntegrationUpdateStatusResponse, error)
	Delete(ctx context.Context, id string) error
	GetByOwnerID(ctx context.Context, req *entity.IntegrationRequest) (*entity.IntegrationGetResponse, error)
	GetTokenUsageList(ctx context.Context, req *entity.IntegrationListRequest) (*entity.IntegrationListResponse, error)
}
