package integration

import (
	"context"
	"sugurta/internal/entity"
	"sugurta/internal/infrastructure/repository"
	"time"
)

type integrationService struct {
	ctxTimeout      time.Duration
	integrationRepo repository.Integration
}

func NewIntegrationService(ctxTimeout time.Duration, repo repository.Integration) repository.Integration {
	return &integrationService{
		ctxTimeout:      ctxTimeout,
		integrationRepo: repo,
	}
}

func (s *integrationService) Create(ctx context.Context, req *entity.IntegrationCreate) error {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	return s.integrationRepo.Create(ctx, req)
}

func (s *integrationService) Update(ctx context.Context, req *entity.IntegrationUpdate) (*entity.IntegrationUpdateResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	return s.integrationRepo.Update(ctx, req)
}

func (s *integrationService) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	return s.integrationRepo.Delete(ctx, id)
}

func (s *integrationService) GetByOwnerID(ctx context.Context, req *entity.IntegrationRequest) (*entity.IntegrationGetResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	return s.integrationRepo.GetByOwnerID(ctx, req)
}

func (s *integrationService) UpdateStatus(ctx context.Context, req *entity.IntegrationUpdateStatus) (*entity.IntegrationUpdateStatusResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	return s.integrationRepo.UpdateStatus(ctx, req)
}


func (s *integrationService) GetTokenUsageList(ctx context.Context, req *entity.IntegrationListRequest) (*entity.IntegrationListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	return s.integrationRepo.GetTokenUsageList(ctx, req)
}

func (s *integrationService) CheckIntegrationExistence(ctx context.Context, businessID string) (*entity.IntegrationExistenceResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	return s.integrationRepo.CheckIntegrationExistence(ctx, businessID)
}


