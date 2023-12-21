package core

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gtngzlv/gophkeeper-server/internal/domain/models"
)

const (
	Secret = "1234"
	userID = "uid"
)

func NewToken(ctx context.Context, user models.User, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims[userID] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	setContextUserID(ctx, user.ID)

	tokenString, err := token.SignedString([]byte(Secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func setContextUserID(ctx context.Context, userID int64) {
	context.WithValue(ctx, userID, userID)
}

func GetContextUserID(ctx context.Context) int64 {
	var id int64
	if value, ok := ctx.Value(userID).(int64); ok {
		id = value
	} else {
		return 0
	}
	return id
}
