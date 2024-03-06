package repository

import (
	"context"
	"log/slog"

	_ "github.com/lib/pq"

	"github.com/gtngzlv/gophkeeper-server/internal/config"
	"github.com/gtngzlv/gophkeeper-server/internal/domain/models"
	"github.com/gtngzlv/gophkeeper-server/internal/repository/postgres"
)

type IRepository interface {
	Login(ctx context.Context, email string) (models.User, error)
	Register(ctx context.Context, email string, passHash []byte, secretKeyHash []byte, encryptedKey []byte) (int64, error)
	SaveData(ctx context.Context, data models.PersonalData, userID int64) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type Repository struct {
	log *slog.Logger
	IRepository
}

func New(ctx context.Context, log *slog.Logger, cfg *config.Config) *Repository {
	db, err := postgres.New(ctx, log, cfg.DBConnectionPath)
	if err != nil {
		log.Error("failed to init db", err.Error())
		return nil
	}
	return &Repository{
		log:         log,
		IRepository: db,
	}
}
