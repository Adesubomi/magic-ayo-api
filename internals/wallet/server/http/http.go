package http

import (
	"fmt"
	walletEntity "github.com/Adesubomi/magic-ayo-api/internals/wallet/entity"
	authPkg "github.com/Adesubomi/magic-ayo-api/pkg/auth"
	configPkg "github.com/Adesubomi/magic-ayo-api/pkg/config"
	"github.com/Adesubomi/magic-ayo-api/pkg/datasource"
	lightningPkg "github.com/Adesubomi/magic-ayo-api/pkg/lightning"
	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	Config      *configPkg.Config
	DbClient    *gorm.DB
	RedisClient *redis.Client
	LndClient   *lightningPkg.LNDClient
}

func (s Service) RegisterRoutes() *fiber.App {
	s.migrateEntities()
	authReg := authPkg.Registry{
		RedisClient: s.RedisClient,
	}

	walletHttp := Handler{
		Config:      s.Config,
		DbClient:    s.DbClient,
		RedisClient: s.RedisClient,
		LndClient:   s.LndClient,
	}

	app := fiber.New()
	app.Get("/", authReg.AuthMiddleware, walletHttp.GetUserWallet)
	app.Post("/make-payment/ln-url", authReg.AuthMiddleware, walletHttp.GenerateInvoiceLNUrl)
	app.Post("/get-status/ln-url/:url", authReg.AuthMiddleware, walletHttp.GetLNInvoiceStatus)
	return app
}

func (s Service) migrateEntities() {
	fmt.Println("")
	fmt.Println("  [...] Migrating tables - wallet")

	entities := []interface{}{
		&walletEntity.Wallet{},
		&walletEntity.BitcoinWallet{},
		&walletEntity.BitcoinAddress{},
		&walletEntity.Transaction{},
	}

	datasource.MigrateTables(s.DbClient, entities)
}
