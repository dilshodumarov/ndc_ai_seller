package settings

import (
	"context"
	"sugurta/internal/entity"
	storage "sugurta/internal/infrastructure/repository"
	"time"
)

type settingsService struct {
	ctxTimeout    time.Duration
	settingsRepo  storage.SettingsStorage
}

func NewSettingsService(ctxTimeout time.Duration, s storage.SettingsStorage) *settingsService {
	return &settingsService{
		ctxTimeout:   ctxTimeout,
		settingsRepo: s,
	}
}

func (ss *settingsService) Create(ctx context.Context, req *entity.CreateOrderStatusRequest) error {
	ctx, cancel := context.WithTimeout(ctx, ss.ctxTimeout)
	defer cancel()

	return ss.settingsRepo.Create(ctx, req)
}

func (ss *settingsService) Get(ctx context.Context, guid string) (*entity.OrderStatus, error) {
	ctx, cancel := context.WithTimeout(ctx, ss.ctxTimeout)
	defer cancel()

	return ss.settingsRepo.Get(ctx, guid)
}

func (ss *settingsService) Update(ctx context.Context, req *entity.UpdateOrderStatusRequest) error {
	ctx, cancel := context.WithTimeout(ctx, ss.ctxTimeout)
	defer cancel()

	return ss.settingsRepo.Update(ctx, req)
}

func (ss *settingsService) Delete(ctx context.Context, guid string) error {
	ctx, cancel := context.WithTimeout(ctx, ss.ctxTimeout)
	defer cancel()

	return ss.settingsRepo.Delete(ctx, guid)
}

func (ss *settingsService) List(ctx context.Context, businessID string) ([]*entity.OrderStatus, error) {
	ctx, cancel := context.WithTimeout(ctx, ss.ctxTimeout)
	defer cancel()

	return ss.settingsRepo.List(ctx, businessID)
}


func (ss *settingsService) GetStatusByName(ctx context.Context, name,bussnesid string) (*string, error) {
	ctx, cancel := context.WithTimeout(ctx, ss.ctxTimeout)
	defer cancel()

	return ss.settingsRepo.GetStatusByName(ctx, name,bussnesid)
}


func (s *settingsService) CreateSettings(ctx context.Context, req *entity.CreateSettingsRequest) error {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	return s.settingsRepo.CreateSettings(ctx, req)
}

func (s *settingsService) GetSettings(ctx context.Context, guid string) (*entity.Settings, error) {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	return s.settingsRepo.GetSettings(ctx, guid)
}

func (s *settingsService) UpdateSettings(ctx context.Context, req *entity.UpdateSettingsRequest) error {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	return s.settingsRepo.UpdateSettings(ctx, req)
}

func (s *settingsService) DeleteSettings(ctx context.Context, guid string) error {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	return s.settingsRepo.DeleteSettings(ctx, guid)
}

func (s *settingsService) ListSettingsByBusinessID(ctx context.Context, businessID string) ([]*entity.Settings, error) {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	return s.settingsRepo.ListSettingsByBusinessID(ctx, businessID)
}

func (s *settingsService)  GetSettingsBussnesId(ctx context.Context, businessID string) (*entity.Settings, error) {
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	return s.settingsRepo.GetSettingsBussnesId(ctx,  businessID)
}