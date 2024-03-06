package grpcapp

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/gtngzlv/gophkeeper-server/internal/config"
	"github.com/gtngzlv/gophkeeper-server/internal/domain/models"
	"github.com/gtngzlv/gophkeeper-server/internal/grpc/gophkeeper"
)

type App struct {
	listener   net.Listener
	log        *slog.Logger
	grpcServer *grpc.Server

	config *config.Config
}

type IGophkeeperService interface {
	Register(ctx context.Context, email string, password string) (userID int64, err error)
	Login(ctx context.Context, email string, password string) (token string, err error)
	SaveData(ctx context.Context, data models.PersonalData) error
}

func New(log *slog.Logger, srv IGophkeeperService, cfg *config.Config) *App {
	grpcServer := grpc.NewServer()

	gophkeeper.Register(grpcServer, srv)

	return &App{
		log:        log,
		grpcServer: grpcServer,
		config:     cfg,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"
	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.config.GRPC.Port))

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.config.GRPC.Port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	reflection.Register(a.grpcServer)
	if err := a.grpcServer.Serve(listener); err != nil {
		log.Error("can't start gRPC server" + err.Error())
		return err
	}
	return nil
}

func (a *App) Stop(ctx context.Context) {
	const op = "grpcapp.Stop"
	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.config.GRPC.Port))
	log.InfoContext(ctx, "stopping gRPC server", a.listener.Addr().String())

	a.grpcServer.GracefulStop()
}
