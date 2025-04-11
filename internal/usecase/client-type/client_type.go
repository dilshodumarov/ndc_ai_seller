package clienttype

import (
	"sugurta/internal/entity"
	"sugurta/internal/infrastructure/repository"
	"context"
	"errors"
	"time"
)

type clientTypeService struct {
	ctxTimeout        time.Duration
	clientTypeService repository.ClientType
}

func NewClientTypeService(ctxTimeout time.Duration, clientType repository.ClientType) ClientType {
	return &clientTypeService{
		ctxTimeout:        ctxTimeout,
		clientTypeService: clientType,
	}
}

// CreateClientType creates a new client type
func (c *clientTypeService) Create(ctx context.Context, clientType *entity.CreateClientTypeRequest) error {
	ctx, cancel := context.WithTimeout(ctx, c.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	if clientType.Name == "" {
		return errors.New("name is required")
	}

	return c.clientTypeService.Create(ctx, clientType)
}

// GetClientType gets a client type by id
func (c *clientTypeService) Get(ctx context.Context, id string) (*entity.ClientType, error) {
	ctx, cancel := context.WithTimeout(ctx, c.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	return c.clientTypeService.Get(ctx, id)
}

func (c *clientTypeService) List(ctx context.Context, limit, page int) ([]*entity.ClientType, error) {
	ctx, cancel := context.WithTimeout(ctx, c.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	return c.clientTypeService.List(ctx, page, limit)
}

// UpdateClientType updates a client type
func (c *clientTypeService) Update(ctx context.Context, clientType *entity.UpdateClientType) error {
	ctx, cancel := context.WithTimeout(ctx, c.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	if clientType.ID == "" {
		return errors.New("id is required")
	}
	if clientType.Name == "" {
		return errors.New("name is required")
	}
	return c.clientTypeService.Update(ctx, clientType)
}

// DeleteClientType deletes a client type
func (c *clientTypeService) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, c.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	return c.clientTypeService.Delete(ctx, id)
}
