package integration

// import (
// 	"context"
// 	"fmt"

// 	"sugurta/internal/entity"
// 	"sugurta/internal/repo"
// 	"sugurta/internal/pkg/logger"
// )

// // UseCase is the auth use case implementation
// type IntegrationUseCase struct {
// 	authRepo        repo.AuthRepo
// 	productRepo     repo.ProductRepo
// 	integrationRepo repo.IntegrationRepo
// 	log             logger.Interface
// }

// // New creates a new integration use case
// func NewIntegrationUseCase(a repo.AuthRepo, p repo.ProductRepo, i repo.IntegrationRepo, l logger.Interface) *IntegrationUseCase {
// 	return &IntegrationUseCase{
// 		authRepo:        a,
// 		productRepo:     p,
// 		integrationRepo: i,
// 		log:             l,
// 	}
// }

// // -------------- Integration --------------
// // CreateIntegration -.
// func (in *IntegrationUseCase) CreateIntegration(ctx context.Context, i entity.Integration) error {
// 	err := in.integration.CreateIntegration(ctx, i)
// 	if err != nil {
// 		return fmt.Errorf("IntegrationUseCase - CreateIntegration - s.integration.CreateIntegration: %w", err)
// 	}

// 	return nil
// }

// // GetIntegration -.
// func (in *IntegrationUseCase) GetIntegration(ctx context.Context, id string) (entity.Integration, error) {
// 	integration, err := uc.integration.GetIntegration(ctx, id)
// 	if err != nil {
// 		return entity.Integration{}, fmt.Errorf("IntegrationUseCase - GetIntegration - s.integration.GetIntegration: %w", err)
// 	}

// 	return integration, nil
// }

// // UpdateIntegration -.
// func (in *IntegrationUseCase) UpdateIntegration(ctx context.Context, i entity.Integration) error {
// 	err := uc.integration.UpdateIntegration(ctx, i)
// 	if err != nil {
// 		return fmt.Errorf("IntegrationUseCase - UpdateIntegration - s.integration.UpdateIntegration: %w", err)
// 	}

// 	return nil
// }

// // DeleteIntegration -.
// func (in *IntegrationUseCase) DeleteIntegration(ctx context.Context, id string) error {
// 	err := uc.integration.DeleteIntegration(ctx, id)
// 	if err != nil {
// 		return fmt.Errorf("IntegrationUseCase - DeleteIntegration - s.integration.DeleteIntegration: %w", err)
// 	}

// 	return nil
// }
