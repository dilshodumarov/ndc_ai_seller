package settings

import (
	"context"
	"sugurta/internal/entity"
)

type SettingsStorage interface {
	Create(ctx context.Context, req *entity.CreateOrderStatusRequest) error
	Get(ctx context.Context, guid string) (*entity.OrderStatus, error)
	Update(ctx context.Context, req *entity.UpdateOrderStatusRequest) error
	Delete(ctx context.Context, guid string) error
	List(ctx context.Context, businessID string) ([]*entity.OrderStatus, error)
	GetStatusByName(ctx context.Context, name,bussnesid string) (*string, error)
	CreateSettings(ctx context.Context, req *entity.CreateSettingsRequest) error
	GetSettings(ctx context.Context, guid string) (*entity.Settings, error)
	UpdateSettings(ctx context.Context, req *entity.UpdateSettingsRequest) error
	DeleteSettings(ctx context.Context, guid string) error
	ListSettingsByBusinessID(ctx context.Context, businessID string) ([]*entity.Settings, error)
	GetSettingsByName(ctx context.Context, name, businessID string) (*entity.Settings, error)
}
