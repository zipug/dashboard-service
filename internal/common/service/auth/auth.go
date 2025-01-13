package auth

import (
	"context"
	"dashboard/internal/common/service/config"
	"errors"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	tokenAuth *jwtauth.JWTAuth
}

func NewValidator() jwt.Validator {
	validator := jwt.ValidatorFunc(func(_ context.Context, t jwt.Token) jwt.ValidationError {
		exp := t.Expiration()
		if exp.Before(time.Now()) {
			err := errors.New("token is expired")
			return jwt.NewValidationError(err)
		}
		return nil
	})
	return validator
}

func New(cfg *config.AppConfig) *Auth {
	token := jwtauth.New(
		"HS256",
		[]byte(cfg.JwtSecretKey),
		nil,
		jwt.WithValidator(NewValidator()),
	)
	return &Auth{tokenAuth: token}
}

func (a *Auth) GetTokenAuth() *jwtauth.JWTAuth {
	return a.tokenAuth
}

func (a *Auth) HashPassword(password string) (string, error) {
	b := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(b, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (a *Auth) ValidatePassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
