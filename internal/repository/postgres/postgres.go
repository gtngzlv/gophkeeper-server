package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pressly/goose"

	customerr "github.com/gtngzlv/gophkeeper-server/internal/domain/errors"
	"github.com/gtngzlv/gophkeeper-server/internal/domain/models"
)

type Postgres struct {
	log *slog.Logger
	db  *sqlx.DB
}

// New creates a new instance of the PostgreSQL storage
func New(ctx context.Context, log *slog.Logger, connString string) (*Postgres, error) {
	const op = "storage.postgres.New"
	log = log.With(
		slog.String("op", op))

	db, err := sqlx.ConnectContext(ctx, "postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	if err = goose.SetDialect("postgres"); err != nil {
		log.Error("unable to set goose dialect", err.Error())
		return nil, err
	}
	if err = goose.Up(db.DB, "migrations"); err != nil {
		log.Error("failed to load migrations ", err.Error())
		return nil, err
	}

	return &Postgres{
		log: log,
		db:  db,
	}, nil
}

func (r *Postgres) Stop() error {
	return r.db.Close()
}

func (r *Postgres) Login(ctx context.Context, email string) (models.User, error) {
	const op = "storage.postgres.Login"
	log := r.log.With(
		slog.String("op", op),
		slog.String("email", email))

	var user models.User
	query := "SELECT ID, EMAIL, PASSWORD_HASH FROM USERS WHERE EMAIL=$1"
	res := r.db.QueryRowContext(ctx, query, email)
	err := res.Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, customerr.ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("%s:%w", op, err)
	}
	log.Info("found user by email")
	return user, nil
}

func (r *Postgres) Register(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "storage.postgres.Register"
	log := r.log.With(
		slog.String("op", op),
		slog.String("email", email))
	var userID int64
	query := "INSERT INTO USERS(email, password_hash) VALUES ($1, $2) RETURNING ID"

	res := r.db.QueryRowContext(ctx, query, email, passHash)
	err := res.Scan(&userID)
	if err != nil {
		if err.(*pq.Error).Code == pgerrcode.UniqueViolation {
			return 0, fmt.Errorf("%s:%w", op, customerr.ErrUserExists)
		}
		return 0, fmt.Errorf("%s:%w", op, err)
	}
	log.Info("registered new user with ID", userID)
	return userID, nil
}

func (r *Postgres) Logout(ctx context.Context) {
	panic("implement me")
}

func (r *Postgres) SaveData(ctx context.Context, data []models.KeeperData) error {
	panic("implement me")
}
