package http

import (
	configPkg "github.com/Adesubomi/magic-ayo-api/pkg/config"
	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Handler struct {
	Config      *configPkg.Config
	DbClient    *gorm.DB
	RedisClient *redis.Client
}

func (h Handler) StartGame(*fiber.Ctx) error {

	return nil
}

func (h Handler) PotPack(*fiber.Ctx) error {

	return nil
}

func (h Handler) AbortGame(*fiber.Ctx) error {

	return nil
}

func (h Handler) GameStatus(ctx *fiber.Ctx) error {
	return nil
}
