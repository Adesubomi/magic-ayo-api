package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"time"
)

type JwtTokenClaim struct {
	jwt.MapClaims
	ID      string `json:"id"`
	Medium  string `json:"medium"`
	MaxTime int64  `json:"maxTime"`
}

func (c JwtTokenClaim) ToMap() jwt.Claims {
	return jwt.MapClaims{
		"id":      c.ID,
		"medium":  c.Medium,
		"maxTime": c.MaxTime,
	}
}

type UserSession struct {
	Token     string    `json:"token"`
	User      User      `json:"user"`
	ExpiresAt time.Time `json:"expiresAt"`
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
		"expiresAt": s.ExpiresAt.String(),
	}
}
