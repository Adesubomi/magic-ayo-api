package http

import (
	"github.com/Adesubomi/magic-ayo-api/internals/auth/handler"
	"github.com/Adesubomi/magic-ayo-api/pkg/config"
	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	Config      *config.Config
	DbClient    *gorm.DB
	RedisClient *redis.Client
}

func (s Service) RegisterRoutes() *fiber.App {
	authHandler := handler.AuthHandler{
		Config:      s.Config,
		DbClient:    s.DbClient,
		RedisClient: s.RedisClient,
	}

	app := fiber.New()
	app.Post("/login", authHandler.Login)
	app.Post("/sign-up", authHandler.SignUp)
	return app
}
