package telegram

import (
	"context"
	"sugurta/internal/entity"
)

type TelegramAccount interface {
	Create(ctx context.Context, req entity.CreateTelegramAccountRequest) (string, error)
	Get(ctx context.Context, id string) (*entity.TelegramAccount, error)
	Update(ctx context.Context, req entity.UpdateTelegramAccountRequest) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, businessID string) ([]entity.TelegramAccount, error)
}
