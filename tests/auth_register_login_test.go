package tests

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gtngzlv/gophkeeper-server/internal/proto/pb"

	libjwt "github.com/gtngzlv/gophkeeper-server/internal/lib/core"
	"github.com/gtngzlv/gophkeeper-server/tests/suite"
)

func TestRegisterLogin_Login_Success(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := fakePassword()

	respRegister, err := st.AuthClient.Register(ctx, &pb.RegisterRequest{
		Email:    email,
		Password: password,
	})

	require.NoError(t, err)

	assert.NotEmpty(t, respRegister.GetUserId())

	respLogin, err := st.AuthClient.Login(ctx, &pb.LoginRequest{
		Email:    email,
		Password: password,
	})

	loginTime := time.Now()

	require.NoError(t, err)
	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(libjwt.Secret), nil
	})
	require.NoError(t, err)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, respRegister.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"].(string))

	const delta = 1

	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), delta)
}

func TestRegisterLogin_LoginWithIncorrectPassword_Failed(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	registerPassword := fakePassword()

	respRegister, err := st.AuthClient.Register(ctx, &pb.RegisterRequest{
		Email:    email,
		Password: registerPassword,
	})

	require.NoError(t, err)

	assert.NotEmpty(t, respRegister.GetUserId())

	loginPassword := fakePassword()
	_, err = st.AuthClient.Login(ctx, &pb.LoginRequest{
		Email:    email,
		Password: loginPassword,
	})
	require.Error(t, err)
}

func TestRegisterLogin_DuplicateRegisterRequest_And_Login(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := fakePassword()

	firstRespRegister, err := st.AuthClient.Register(ctx, &pb.RegisterRequest{
		Email:    email,
		Password: password,
	})

	require.NoError(t, err)

	assert.NotEmpty(t, firstRespRegister.GetUserId())

	_, err = st.AuthClient.Register(ctx, &pb.RegisterRequest{
		Email:    email,
		Password: password,
	})

	require.Error(t, err)

	respLogin, err := st.AuthClient.Login(ctx, &pb.LoginRequest{
		Email:    email,
		Password: password,
	})

	loginTime := time.Now()

	require.NoError(t, err)
	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(libjwt.Secret), nil
	})
	require.NoError(t, err)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, firstRespRegister.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"].(string))

	const delta = 1

	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), delta)
}

func fakePassword() string {
	passwordLength := 8
	password := gofakeit.Password(true, true, true, true, true, passwordLength)
	return password
}
