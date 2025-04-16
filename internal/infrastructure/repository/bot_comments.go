package repository

import (
	"context"
	"sugurta/internal/entity"
)

type BotCommandStorage interface {
	CreateBotCommand(ctx context.Context, cmd entity.BotCommandRequest) (string, error)
	GetBotCommand(ctx context.Context, guid string) (*entity.BotCommandResponse, error)
	UpdateBotCommand(ctx context.Context, cmd entity.BotCommandUpdateRequest) error
	DeleteBotCommand(ctx context.Context, guid string) error
	ListBotCommands(ctx context.Context, integrationID string) ([]entity.BotCommandResponse, error)
}
