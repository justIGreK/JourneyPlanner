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
		Expiration: time.Now().Add(HoursInDay * time.Hour),
	}

	encrypted, err := pasetoInstance.Encrypt(symmetricKey, payload, nil)
	if err != nil{
		logs.Error(err)
		return "", fmt.Errorf("GeneratePasetoToken error: %v", err)
	}
	return encrypted, nil
}

func (s *UserSrv) ValidatePasetoToken(tokenString string) (*TokenPayload, error) {
	symmetricKey := []byte(os.Getenv("SYMMETRIC_KEY"))
	var payload TokenPayload
	var footer string
	err := pasetoInstance.Decrypt(tokenString, symmetricKey, &payload, &footer)
	if err != nil {
		logs.Error(err)
		return nil, fmt.Errorf("validatePasetoToken error: %v", err)
	}

	if time.Now().After(payload.Expiration) {
		return nil, fmt.Errorf("token expired")
	}
	return &payload, nil
}
