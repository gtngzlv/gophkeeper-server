package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gtngzlv/gophkeeper-server/internal/domain/models"
)

const secret = "1234"

func NewToken(user models.User, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
