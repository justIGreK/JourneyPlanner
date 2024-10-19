package service

import (
	"fmt"
	"os"
	"time"

	"github.com/o1egl/paseto"
)

var pasetoInstance = paseto.NewV2()

type TokenPayload struct {
	UserLogin  string    `json:"user_login"`
	Expiration time.Time `json:"expiration"`
}

func (s *UserSrv) GeneratePasetoToken(userLogin string) (string, error) {
	symmetricKey := []byte(os.Getenv("SYMMETRIC_KEY"))
	payload := TokenPayload{
		UserLogin:  userLogin,
		Expiration: time.Now().Add(24 * time.Hour),
	}

	encrypted, err := pasetoInstance.Encrypt(symmetricKey, payload, nil)
	return encrypted, err
}

func ValidatePasetoToken(tokenString string) (*TokenPayload, error) {
	symmetricKey := []byte(os.Getenv("SYMMETRIC_KEY"))
	var payload TokenPayload
	var footer string
	err := pasetoInstance.Decrypt(tokenString, symmetricKey, &payload, &footer)
	if err != nil {
		return nil, err
	}

	if time.Now().After(payload.Expiration) {
		return nil, fmt.Errorf("token expired")
	}
	return &payload, nil
}
