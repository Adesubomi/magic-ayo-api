package http

import (
	authPkg "github.com/Adesubomi/magic-ayo-api/pkg/auth"
	configPkg "github.com/Adesubomi/magic-ayo-api/pkg/config"
	responsePkg "github.com/Adesubomi/magic-ayo-api/pkg/response"
	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Handler struct {
	Config      *configPkg.Config
	DbClient    *gorm.DB
	RedisClient *redis.Client
}

func (h Handler) getWalletRepo(ctx *fiber.Ctx) {
	return
}

func (h Handler) GetUserWallet(ctx *fiber.Ctx) error {
	userSession := authPkg.UserSessionFromFiberCtx(ctx)
	if userSession == nil {
		return responsePkg.Unauthorized(ctx, "user information not found")
	}

	// get user wallet

	return responsePkg.Success(
		ctx,
		"user wallet",
		map[string]interface{}{
			"userSession": userSession.SerializeForResponse(),
		},
	)
}
