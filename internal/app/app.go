package app

import (
	"context"
	"log/slog"

	grpcapp "github.com/gtngzlv/gophkeeper-server/internal/app/grpc"
	"github.com/gtngzlv/gophkeeper-server/internal/config"
	"github.com/gtngzlv/gophkeeper-server/internal/repository"
	"github.com/gtngzlv/gophkeeper-server/internal/services/gophkeeper"
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

	srv := gophkeeper.New(log, repo, cfg.TokenTTL)
	grpcApp := grpcapp.New(log, srv, cfg)

	return &App{
		GRPCSrv: grpcApp,
	}, nil
}
