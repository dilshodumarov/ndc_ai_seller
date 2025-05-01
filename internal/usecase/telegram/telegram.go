package telegram

import (
	"context"
	"sugurta/internal/entity"
	"sugurta/internal/infrastructure/repository"
	"time"
)

type telegramService struct {
	ctxTimeout     time.Duration
	telegramRepo   repository.TelegramAccount
}

func NewTelegramService(ctxTimeout time.Duration, r repository.TelegramAccount) *telegramService {
	return &telegramService{
		ctxTimeout:   ctxTimeout,
		telegramRepo: r,
	}
}

func (s *telegramService) Create(ctx context.Context, req entity.CreateTelegramAccountRequest) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	return s.telegramRepo.Create(ctx, req)
}

func (s *telegramService) Get(ctx context.Context, id string) (*entity.TelegramAccount, error) {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	return s.telegramRepo.Get(ctx, id)
}

func (s *telegramService) Update(ctx context.Context, req entity.UpdateTelegramAccountRequest) error {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	return s.telegramRepo.Update(ctx, req)
}

func (s *telegramService) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	return s.telegramRepo.Delete(ctx, id)
}

func (s *telegramService) List(ctx context.Context, businessID string) ([]entity.TelegramAccount, error) {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	return s.telegramRepo.List(ctx, businessID)
}
