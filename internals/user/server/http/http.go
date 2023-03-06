package http

import (
	"fmt"
	"github.com/Adesubomi/magic-ayo-api/pkg/auth"
	"github.com/Adesubomi/magic-ayo-api/pkg/config"
	"github.com/Adesubomi/magic-ayo-api/pkg/datasource"
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
	s.migrateEntities()

	authHandler := Handler{
		Config:      s.Config,
		DbClient:    s.DbClient,
		RedisClient: s.RedisClient,
	}

	app := fiber.New()
	app.Post("/login", authHandler.Login)
	app.Post("/sign-up", authHandler.SignUp)
	return app
}

func (s Service) migrateEntities() {
	fmt.Println("")
	fmt.Println("  [...] Migrating tables - auth")
	entities := []interface{}{
		&auth.User{},
	}

	datasource.MigrateTables(
		s.DbClient,
		entities)
}
