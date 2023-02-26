package auth

import (
	"github.com/Adesubomi/magic-ayo-api/pkg/response"
	"github.com/Adesubomi/magic-ayo-api/pkg/util"
	"github.com/gofiber/fiber/v2"
)

func (r Registry) AuthMiddleware(ctx *fiber.Ctx) error {
	bearerToken := util.GetBearerTokenFromAuthorizationHeader(ctx)

	userSession, err := r.getUserAuthSession(bearerToken)
	if err != nil || userSession == nil || userSession.User.ID == "" {
		return response.Unauthorized(ctx, "Service is only available for logged in users")
	}

	ctx.Locals("BearerToken", bearerToken)
	ctx.Locals("UserSession", userSession)

	return ctx.Next()
}

func (r Registry) GuestMiddleware(ctx *fiber.Ctx) error {
	bearerToken := util.GetBearerTokenFromAuthorizationHeader(ctx)
	userSession, err := r.getUserAuthSession(bearerToken)

	if err == nil {
		return response.Unauthorized(ctx, "Service is only available to guests")
	}

	if userSession != nil && userSession.User.ID != "" {
		return response.Unauthorized(ctx, "Service is only available to guests")
	}

	return ctx.Next()
}
