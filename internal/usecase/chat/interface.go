package chat

import (
	"context"
	"sugurta/internal/entity"
)

type Chat interface {
	Create(ctx context.Context, h *entity.ChatHistory) error
	Get(ctx context.Context, guid string) (*entity.ChatHistory, error)
	List(ctx context.Context, chatID int64) ([]*entity.ChatHistory, error)
	Update(ctx context.Context, h *entity.ChatHistory) error
	Delete(ctx context.Context, guid string) error
}
