package authservice

import (
	"context"
	"log/slog"
)

type IAuthService interface {
	Register(ctx context.Context, email string, password string) (userID int64, err error)
	Login(ctx context.Context, email string, password string) (token string, err error)
	Logout(ctx context.Context)
}

type Service struct {
	auth   IAuthService
	logger *slog.Logger
}

func New(logger *slog.Logger, authService IAuthService) *Service {
	return &Service{
		auth:   authService,
		logger: logger,
	}
}

func (s *Service) Register(ctx context.Context, email string, password string) (userID int64, err error) {
}

func (s *Service) Login(ctx context.Context, email string, password string) (token string, err error) {
}

func (s *Service) Logout(ctx context.Context) {}
