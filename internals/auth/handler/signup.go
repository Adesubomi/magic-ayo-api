package handler

import (
	responsePkg "github.com/Adesubomi/magic-ayo-api/pkg/response"
	"github.com/gofiber/fiber/v2"
)

func (h AuthHandler) SignUp(ctx *fiber.Ctx) error {
	return responsePkg.Success(
		ctx,
		"Sign up successful.",
		map[string]interface{}{})
}
