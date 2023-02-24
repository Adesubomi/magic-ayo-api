package cmd

import (
	"fmt"
	authHttp "github.com/Adesubomi/magic-ayo-api/internals/auth/server/http"
	walletHttp "github.com/Adesubomi/magic-ayo-api/internals/wallet/server/http"
	pkgConfig "github.com/Adesubomi/magic-ayo-api/pkg/config"
	dataPkg "github.com/Adesubomi/magic-ayo-api/pkg/datasource"
	logPkg "github.com/Adesubomi/magic-ayo-api/pkg/log"
	"github.com/BurntSushi/toml"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
)

func main() {
	var config *pkgConfig.Config
	if _, err := toml.DecodeFile("config.toml", config); err != nil {
		fmt.Println("Error decoding TOML file:", err)
		return
	}

	dbClient, err := dataPkg.ConnectDatabase(config)
	if err != nil {
		logPkg.ReportError(err)
		os.Exit(-1)
	}

	redisClient, err := dataPkg.RedisConnection(config)
	if err != nil {
		logPkg.ReportError(err)
		os.Exit(-1)
	}

	auth := authHttp.Service{
		Config:      config,
		DbClient:    dbClient,
		RedisClient: redisClient,
	}

	wallet := walletHttp.Service{
		Config:      config,
		DbClient:    dbClient,
		RedisClient: redisClient,
	}

	app := fiber.New()
	app.Mount("auth", auth.RegisterRoutes())
	app.Mount("wallets", wallet.RegisterRoutes())

	err = app.Listen(fmt.Sprintf(":%v", config.AppPort))
	if err != nil {
		log.Fatalf("server connection listen and server failed, error message %s", err.Error())
	}
}
