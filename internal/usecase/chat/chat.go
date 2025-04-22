package chat

import (
	"context"
	"time"

	"sugurta/internal/entity"
	"sugurta/internal/infrastructure/repository"
)

type chatService struct {
	ctxTimeout time.Duration
	chatRepo   repository.Chat
}

func NewChatService(ctxTimeout time.Duration, repo repository.Chat) Chat {
	return &chatService{
		ctxTimeout: ctxTimeout,
		chatRepo:   repo,
	}
}

func (c *chatService) Create(ctx context.Context, h *entity.ChatHistory) error {
	ctx, cancel := context.WithTimeout(ctx, c.ctxTimeout)
	defer cancel()

	return c.chatRepo.Create(ctx, h)
}

func (c *chatService) Get(ctx context.Context, guid string) (*entity.ChatHistory, error) {
	ctx, cancel := context.WithTimeout(ctx, c.ctxTimeout)
	defer cancel()

	return c.chatRepo.Get(ctx, guid)
}

func (c *chatService) List(ctx context.Context, chatID int64) ([]*entity.ChatHistory, error) {
	ctx, cancel := context.WithTimeout(ctx, c.ctxTimeout)
	defer cancel()

	return c.chatRepo.List(ctx, chatID)
}

func (c *chatService) Update(ctx context.Context, h *entity.ChatHistory) error {
	ctx, cancel := context.WithTimeout(ctx, c.ctxTimeout)
	defer cancel()

	return c.chatRepo.Update(ctx, h)
}

func (c *chatService) Delete(ctx context.Context, guid string) error {
	ctx, cancel := context.WithTimeout(ctx, c.ctxTimeout)
	defer cancel()

	return c.chatRepo.Delete(ctx, guid)
}
