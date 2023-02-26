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

	gameHttp := Handler{
		Config:      s.Config,
		DbClient:    s.DbClient,
		RedisClient: s.RedisClient,
	}

	app := fiber.New()
	app.Get("/start", authReg.AuthMiddleware, gameHttp.StartGame)
	app.Get("/pot-pack", authReg.AuthMiddleware, gameHttp.PotPack)
	app.Get("/abort", authReg.AuthMiddleware, gameHttp.AbortGame)
	app.Get("/status", authReg.AuthMiddleware, gameHttp.GameStatus)
	return app
}
