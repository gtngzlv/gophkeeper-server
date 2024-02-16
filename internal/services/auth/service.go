package authservice

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
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
	Register(ctx context.Context, email string, passHash []byte, secretKeyHash []byte, encryptedKey []byte) (int64, error)
	Login(ctx context.Context, email string) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
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

	// Генерация секретного ключа
	secretKey, err := generateSecretKey()
	if err != nil {
		log.Error("failed to generate secret key", err.Error())
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	// Хеширование секретного ключа
	secretKeyHash := hashSecretKey(secretKey)

	// Хеширование пароля
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", err.Error())
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	// Шифрование секретного ключа на основе пароля
	encryptedKey, err := encryptSecretKey(secretKey, []byte(password))
	if err != nil {
		log.Error("failed to encrypt secret key", err.Error())
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	userID, err := s.storageAuth.Register(ctx, email, passHash, []byte(secretKeyHash), encryptedKey)
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

	user, err := s.storageAuth.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, customerr.ErrUserNotFound) {
			s.logger.Warn("user not found", err.Error())
			return "", fmt.Errorf("%s:%w", op, customerr.ErrInvalidCredentials)
		}

		log.Error("failed to get user", err.Error())
		return "", fmt.Errorf("%s:%w", op, err)
	}

	decryptedKey, err := decryptSecretKey(user.EncryptedKey, []byte(password))
	if err != nil {
		log.Error("failed to decrypt secret key", err.Error())
		return "", fmt.Errorf("%s:%w", op, customerr.ErrInvalidCredentials)
	}

	if !compareHashes([]byte(user.SecretKeyHash), []byte(hashSecretKey(decryptedKey))) {
		log.Info("invalid credentials")
		return "", fmt.Errorf("%s:%w", op, customerr.ErrInvalidCredentials)
	}

	// Проверка пароля
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Info("invalid credentials", err.Error())
		return "", fmt.Errorf("%s:%w", op, customerr.ErrInvalidCredentials)
	}

	log.Info("user logged in successfuly")

	// Генерация токена
	token, err := core.NewToken(ctx, user, s.tokenTTL)
	if err != nil {
		log.Error("failed to generate token", err.Error())
		return "", fmt.Errorf("%s:%w", op, err)
	}

	return token, nil
}

func generateSecretKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}

func hashSecretKey(key []byte) string {
	hash := sha256.Sum256(key)
	return string(hash[:])
}

func encryptSecretKey(key []byte, password []byte) ([]byte, error) {
	// Преобразование пароля в ключ с использованием хеш-функции
	hashedPassword := sha256.Sum256(password)

	block, err := aes.NewCipher(hashedPassword[:])
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, key, nil)
	ciphertext = append(nonce, ciphertext...)

	return ciphertext, nil
}

// decryptSecretKey расшифровывает секретный ключ на основе пароля.
func decryptSecretKey(ciphertext []byte, password []byte) ([]byte, error) {
	// Преобразование пароля в ключ с использованием хеш-функции
	hashedPassword := sha256.Sum256(password)

	block, err := aes.NewCipher(hashedPassword[:])
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	nonce := ciphertext[:gcm.NonceSize()]
	ciphertext = ciphertext[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// compareHashes сравнивает два хеша без раскрывания конкретного значения.
func compareHashes(hash1, hash2 []byte) bool {
	return subtle.ConstantTimeCompare(hash1, hash2) == 1
}
