package botcomments

import (
	"context"
	"sugurta/internal/entity"
	"sugurta/internal/infrastructure/repository"
	"time"
)

// UseCase is the auth use case implementation
type botCommentsService struct {
	ctxTimeout      time.Duration
	botCommentsRepo repository.BotCommandStorage
}

func NewbotCommentsService(ctxTimeout time.Duration, b repository.BotCommandStorage) *botCommentsService {
	return &botCommentsService{
		ctxTimeout:      ctxTimeout,
		botCommentsRepo: b,
	}
}

func (s *botCommentsService) CreateBotCommand(ctx context.Context, cmd entity.BotCommandRequest) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	return s.botCommentsRepo.CreateBotCommand(ctx, cmd)
}

func (s *botCommentsService) GetBotCommand(ctx context.Context, guid string) (*entity.BotCommandResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	return s.botCommentsRepo.GetBotCommand(ctx, guid)
}

func (s *botCommentsService) UpdateBotCommand(ctx context.Context, cmd entity.BotCommandUpdateRequest) error {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	return s.botCommentsRepo.UpdateBotCommand(ctx, cmd)
}

func (s *botCommentsService) DeleteBotCommand(ctx context.Context, guid string) error {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	return s.botCommentsRepo.DeleteBotCommand(ctx, guid)
}

func (s *botCommentsService) ListBotCommands(ctx context.Context, integrationID string) ([]entity.BotCommandResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	return s.botCommentsRepo.ListBotCommands(ctx, integrationID)
}
