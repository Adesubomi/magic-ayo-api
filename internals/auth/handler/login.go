package handler

import (
	responsePkg "github.com/Adesubomi/magic-ayo-api/pkg/response"
	"github.com/gofiber/fiber/v2"
)

func (h AuthHandler) Login(ctx *fiber.Ctx) error {
	return responsePkg.Success(
		ctx,
		"Login successful",
		map[string]interface{}{})
}
