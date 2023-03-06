package http

import (
	"fmt"
	aayoEntity "github.com/Adesubomi/magic-ayo-api/internals/aayo/entity"
	authPkg "github.com/Adesubomi/magic-ayo-api/pkg/auth"
	configPkg "github.com/Adesubomi/magic-ayo-api/pkg/config"
	"github.com/Adesubomi/magic-ayo-api/pkg/datasource"
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
	s.migrateEntities()

	authReg := authPkg.Registry{
		RedisClient: s.RedisClient,
	}

	gameHttp := Handler{
		Config:      s.Config,
		DbClient:    s.DbClient,
		RedisClient: s.RedisClient,
	}

	app := fiber.New()
	app.Post("/start", authReg.AuthMiddleware, gameHttp.StartGame)
	app.Post("/pot-pack", authReg.AuthMiddleware, gameHttp.PotPack)
	app.Post("/abort", authReg.AuthMiddleware, gameHttp.AbortGame)
	app.Get("/status", authReg.AuthMiddleware, gameHttp.GameStatus)
	return app
}

func (s Service) migrateEntities() {
	fmt.Println("")
	fmt.Println("  [...] Migrating tables - aayo")

	entities := []interface{}{
		&aayoEntity.GamePlay{},
	}

	datasource.MigrateTables(s.DbClient, entities)
}
