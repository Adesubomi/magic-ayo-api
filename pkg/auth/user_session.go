package auth

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

type JwtTokenClaim struct {
	jwt.MapClaims
	ID        string `json:"id"`
	ExpiresAt int64  `json:"expiresAt"`
}

func (c JwtTokenClaim) ToMap() jwt.Claims {
	return jwt.MapClaims{
		"id":        c.ID,
		"expiresAt": c.ExpiresAt,
	}
}

type UserSession struct {
	Token     string `json:"token"`
	User      User   `json:"user"`
	ExpiresAt int64  `json:"expiresAt"`
}

func UserSessionFromFiberCtx(ctx *fiber.Ctx) *UserSession {
	userSession, ok := ctx.Locals("UserSession").(*UserSession)
	if !ok {
		return nil
	}
	return userSession
}

func (s UserSession) SerializeForResponse() map[string]string {
	return map[string]string{
		"token":     s.Token,
		"expiresAt": fmt.Sprintf("%v", s.ExpiresAt),
	}
}

func (s UserSession) IsActive() bool {
	return false
}
