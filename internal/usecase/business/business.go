package business

import (
	"sugurta/internal/entity"
	"sugurta/internal/infrastructure/repository"
	"context"
	"time"
)

// UseCase is the auth use case implementation
type businessService struct {
	ctxTimeout   time.Duration
	businessRepo repository.Business
}

func NewbusinessService(ctxTimeout time.Duration, b repository.Business) Business {
	return &businessService{
		ctxTimeout:   ctxTimeout,
		businessRepo: b,
	}
}

func (bu *businessService) Create(ctx context.Context, b *entity.CreateBusinessRequest) error {
	ctx, cancel := context.WithTimeout(ctx, bu.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	return bu.businessRepo.Create(ctx, b)
}

func (bu *businessService) Get(ctx context.Context, id string) (*entity.Business, error) {
	ctx, cancel := context.WithTimeout(ctx, bu.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	return bu.businessRepo.Get(ctx, id)
}

func (bu *businessService) Update(ctx context.Context, b *entity.UpdateBusinessRequest) error {
	ctx, cancel := context.WithTimeout(ctx, bu.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	return bu.businessRepo.Update(ctx, b)
}

func (bu *businessService) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, bu.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	return bu.businessRepo.Delete(ctx, id)
}

func (bu *businessService) List(ctx context.Context,  busness entity.GetAllBusinessesRequest) (*entity.GetAllBusinessesResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, bu.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	return bu.businessRepo.List(ctx, busness)
}
