package models

type User struct {
	ID            int64
	Email         string
	PassHash      []byte
	SecretKeyHash string
	EncryptedKey  []byte
}
