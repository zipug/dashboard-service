package ports

import "github.com/go-chi/jwtauth/v5"

type Auth interface {
	GetTokenAuth() *jwtauth.JWTAuth
	HashPassword(string) (string, error)
	ValuidatePassword(string, string) bool
}
