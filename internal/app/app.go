package app

import (
	"context"
	"log/slog"

	grpcapp "github.com/gtngzlv/gophkeeper-server/internal/app/grpc"
	"github.com/gtngzlv/gophkeeper-server/internal/config"
	"github.com/gtngzlv/gophkeeper-server/internal/repository/sqlite"
	authservice "github.com/gtngzlv/gophkeeper-server/internal/services/auth"
	keeperservice "github.com/gtngzlv/gophkeeper-server/internal/services/keeper"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func NewApp(
	ctx context.Context,
	log *slog.Logger,
	cfg *config.Config) (*App, error) {
	const op = "App.New"
	log = log.With(slog.String("op", op))

	repo, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init sqlite", err)
		return nil, err
	}

	authSrv := authservice.New(log, repo)
	keeperSrv := keeperservice.New(log, repo)
	grpcApp := grpcapp.New(log, grpcapp.Params{
		AuthService:   authSrv,
		KeeperService: keeperSrv,
	})

	return &App{
		GRPCSrv: grpcApp,
	}, nil
}
