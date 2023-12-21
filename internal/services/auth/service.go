package authservice

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"

	customerr "github.com/gtngzlv/gophkeeper-server/internal/domain/errors"
	"github.com/gtngzlv/gophkeeper-server/internal/domain/models"
	"github.com/gtngzlv/gophkeeper-server/internal/lib/core"
)

type IStorageAuth interface {
	Register(ctx context.Context, email string, passHash []byte) (userID int64, err error)
	Login(ctx context.Context, email string) (models.User, error)
}

type Service struct {
	storageAuth IStorageAuth
	logger      *slog.Logger
	tokenTTL    time.Duration
}

// New returns a new instance of the Auth service
func New(logger *slog.Logger, authService IStorageAuth, tokenTTL time.Duration) *Service {
	return &Service{
		storageAuth: authService,
		logger:      logger,
		tokenTTL:    tokenTTL,
	}
}

// Register creates new user in the system, if email is not exist already. Returns errors, if exists, userID if not.
func (s *Service) Register(ctx context.Context, email string, password string) (int64, error) {
	const op = "service.Auth.Register"

	log := s.logger.With(
		slog.String("op", op),
		slog.String("email", email))

	log.Info("registering new user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", err.Error())

		return 0, fmt.Errorf("%s:%w", op, err)
	}

	userID, err := s.storageAuth.Register(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, customerr.ErrUserExists) {
			log.Warn("user already exists", err.Error())
			return 0, fmt.Errorf("%s:%w", op, customerr.ErrUserExists)
		}
		log.Error("failed to register user", err.Error())
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	log.Info("user registered")
	return userID, nil
}

// Login checks if provided credentials exists in the system and returns token, if yes. Error, if not.
func (s *Service) Login(ctx context.Context, email string, password string) (string, error) {
	const op = "service.Auth.Login"

	log := s.logger.With(
		slog.String("op", op),
		slog.String("email", email))

	log.Info("login")

	user, err := s.storageAuth.Login(ctx, email)
	if err != nil {
		if errors.Is(err, customerr.ErrUserNotFound) {
			s.logger.Warn("user not found", err.Error())
			return "", fmt.Errorf("%s:%w", op, customerr.ErrInvalidCredentials)
		}

		log.Error("failed to get user", err.Error())
		return "", fmt.Errorf("%s:%w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Info("invalid credentials", err.Error())
		return "", fmt.Errorf("%s:%w", op, customerr.ErrInvalidCredentials)
	}

	log.Info("user logged in successfuly")

	token, err := core.NewToken(ctx, user, s.tokenTTL)
	if err != nil {
		log.Error("failed to generate token", err.Error())
		return "", fmt.Errorf("%s:%w", op, err)
	}
	return token, nil
}
