package app

import (
	"context"
	"log/slog"

	grpcapp "github.com/gtngzlv/gophkeeper-server/internal/app/grpc"
	"github.com/gtngzlv/gophkeeper-server/internal/config"
	"github.com/gtngzlv/gophkeeper-server/internal/repository"
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

	repo := repository.New(ctx, log, cfg)

	authSrv := authservice.New(log, repo, cfg.TokenTTL)
	keeperSrv := keeperservice.New(log, repo)
	grpcApp := grpcapp.New(log, grpcapp.Params{
		AuthService:   authSrv,
		KeeperService: keeperSrv,
	}, cfg)

	return &App{
		GRPCSrv: grpcApp,
	}, nil
}
