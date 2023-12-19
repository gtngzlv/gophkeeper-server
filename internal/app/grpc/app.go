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
	"github.com/gtngzlv/gophkeeper-server/internal/grpc/auth"
	"github.com/gtngzlv/gophkeeper-server/internal/grpc/keeper"
)

type App struct {
	listener   net.Listener
	log        *slog.Logger
	grpcServer *grpc.Server

	config config.Config
}

type IAuthService interface {
	Register(ctx context.Context, email string, password string) (userID int64, err error)
	Login(ctx context.Context, email string, password string) (token string, err error)
	Logout(ctx context.Context)
}

type IKeeperService interface {
	SaveData(ctx context.Context, data []models.KeeperData) error
}

type Params struct {
	AuthService   IAuthService
	KeeperService IKeeperService
}

func New(log *slog.Logger, params Params) *App {
	grpcServer := grpc.NewServer()

	auth.Register(grpcServer, params.AuthService)

	keeper.Register(grpcServer, params.KeeperService)

	return &App{
		log:        log,
		grpcServer: grpcServer,
	}
}

func (a *App) Start(ctx context.Context) error {
	const op = "grpcapp.Start"
	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.config.GRPC.Port))
	log.InfoContext(ctx, "starting gRPC server")

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.config.GRPC.Port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	a.listener = listener
	return nil
}

func (a *App) MustRun(ctx context.Context) {
	if err := a.Run(ctx); err != nil {
		panic(err)
	}
}

func (a *App) Run(ctx context.Context) error {
	const op = "grpcapp.Run"
	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.config.GRPC.Port))
	log.InfoContext(ctx, "running gRPC server", a.listener.Addr().String())

	reflection.Register(a.grpcServer)
	if err := a.grpcServer.Serve(a.listener); err != nil {
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
