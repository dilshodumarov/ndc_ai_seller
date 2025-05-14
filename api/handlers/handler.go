package handlers

import (
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
	"golang.org/x/net/context"

	"sugurta/api/middleware"
	"sugurta/internal/infrastructure/repository/redis"
	"sugurta/internal/pkg/config"
	botcomments "sugurta/internal/usecase/bot_comments"
	"sugurta/internal/usecase/business"
	"sugurta/internal/usecase/category"
	"sugurta/internal/usecase/chat"
	clienttype "sugurta/internal/usecase/client-type"
	"sugurta/internal/usecase/intgration"
	"sugurta/internal/usecase/notification"
	"sugurta/internal/usecase/order"
	"sugurta/internal/usecase/product"
	"sugurta/internal/usecase/role"
	"sugurta/internal/usecase/settings"
	"sugurta/internal/usecase/telegram"
	"sugurta/internal/usecase/user"
)

const (
	InvestorToken = "investor"
)

type HandlerOption struct {
	Config         *config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	Enforcer       *casbin.CachedEnforcer
	Cache          redis.Cache
	// Service        grpcClients.ServiceClient
	// RefreshToken   refresh_token.RefreshToken

	User       user.User
	Role       role.Role
	ClientType clienttype.ClientType
	Business   business.Business
	Product    product.Product
	Category   category.Category
	Order      order.Order
	Integration integration.Integration
	BotComments botcomments.BotCommandStorage
	Chat        chat.Chat
	Minio       *minio.Client
	Telegram    telegram.TelegramAccount
	Notification notification.Notification
	Settings     settings.SettingsStorage
}

type BaseHandler struct {
	Cache  redis.Cache
	Config *config.Config
	// Client grpcClients.ServiceClient
}

func (h *BaseHandler) GetAuthData(ctx context.Context) (map[string]string, bool) {
	// tracing
	// ctx, span := otlp.Start(ctx, "handler", "GetAuthData")
	// defer span.End()

	data, ok := ctx.Value(middleware.RequestAuthCtx).(map[string]string)
	return data, ok
}
