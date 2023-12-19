package keeperservice

import (
	"context"
	"log/slog"

	"github.com/gtngzlv/gophkeeper-server/internal/domain/models"
)

type IKeeperService interface {
	SaveData(ctx context.Context, data []models.KeeperData) error
}

type Service struct {
	logger *slog.Logger

	keeper IKeeperService
}

func New(logger *slog.Logger, keeperService IKeeperService) *Service {
	return &Service{
		logger: logger,
		keeper: keeperService,
	}
}

func (s *Service) SaveData(ctx context.Context, data []models.KeeperData) error {}
