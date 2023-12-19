package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gtngzlv/gophkeeper-server/internal/domain/models"
)

type Repository struct {
	db *sql.DB
}

func New(storagePath string) (*Repository, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return &Repository{db: db}, nil
}

func (s *Repository) Stop() error {
	return s.db.Close()
}

func (r *Repository) Login(ctx context.Context, email string, password string) (token string, err error) {

}

func (r *Repository) Register(ctx context.Context, email string, password string) (userID int64, err error) {

}

func (r *Repository) Logout(ctx context.Context) {

}

func (r *Repository) SaveData(ctx context.Context, data []models.KeeperData) error {

}
