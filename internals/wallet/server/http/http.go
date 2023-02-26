package http

import (
	authPkg "github.com/Adesubomi/magic-ayo-api/pkg/auth"
	configPkg "github.com/Adesubomi/magic-ayo-api/pkg/config"
	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	Config      *configPkg.Config
	DbClient    *gorm.DB
	RedisClient *redis.Client
}

func (s Service) RegisterRoutes() *fiber.App {
	authReg := authPkg.Registry{
		RedisClient: s.RedisClient,
	}

	walletHttp := Handler{
		Config:      s.Config,
		DbClient:    s.DbClient,
		RedisClient: s.RedisClient,
	}

	app := fiber.New()
	app.Get("/", authReg.AuthMiddleware, walletHttp.GetUserWallet)
	return app
}
