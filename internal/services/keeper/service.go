package keeperservice

import (
	"context"
	"log/slog"

	customerr "github.com/gtngzlv/gophkeeper-server/internal/domain/errors"
	"github.com/gtngzlv/gophkeeper-server/internal/domain/models"
	"github.com/gtngzlv/gophkeeper-server/internal/lib/core"
)

type IKeeperService interface {
	SaveData(ctx context.Context, data models.PersonalData, userID int64) error
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

func (s *Service) SaveData(ctx context.Context, data models.PersonalData) error {
	const op = "service.Keeper.Register"

	log := s.logger.With(
		slog.String("op", op))

	userID := core.GetContextUserID(ctx)
	if userID == 0 {
		return customerr.ErrFailedGetUserID
	}

	err := s.keeper.SaveData(ctx, data, userID)
	if err != nil {
		log.Error("failed to save data", err.Error())
		return err
	}
	return nil
}
