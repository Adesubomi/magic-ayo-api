package main

import (
	"fmt"
	aayoHttp "github.com/Adesubomi/magic-ayo-api/internals/aayo/server/http"
	authHttp "github.com/Adesubomi/magic-ayo-api/internals/user/server/http"
	walletHttp "github.com/Adesubomi/magic-ayo-api/internals/wallet/server/http"
	pkgConfig "github.com/Adesubomi/magic-ayo-api/pkg/config"
	dataPkg "github.com/Adesubomi/magic-ayo-api/pkg/datasource"
	lnPkg "github.com/Adesubomi/magic-ayo-api/pkg/lightning"
	logPkg "github.com/Adesubomi/magic-ayo-api/pkg/log"
	networkPkg "github.com/Adesubomi/magic-ayo-api/pkg/network"
	"github.com/BurntSushi/toml"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
)

func main() {
	config := &pkgConfig.Config{}
	if _, err := toml.DecodeFile("cmd/config.toml", config); err != nil {
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

	lndSendClient, err := lnPkg.NewSenderLnClient(config)
	if err != nil {
		logPkg.ReportError(err)
		os.Exit(-1)
	}
	if lndSendClient.Connection != nil {
		defer func(lc *lnPkg.LNClient) {
			err := lc.Connection.Close()
			if err != nil {
				logPkg.ReportError(err)
			}
		}(lndSendClient)
	}

	lndReceiveClient, err := lnPkg.NewRecipientLnClient(config)
	if err != nil {
		logPkg.ReportError(err)
		os.Exit(-1)
	}
	if lndReceiveClient.Connection != nil {
		defer func(lc *lnPkg.LNClient) {
			err := lc.Connection.Close()
			if err != nil {
				logPkg.ReportError(err)
			}
		}(lndReceiveClient)
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
		LNClients: &lnPkg.LNClients{
			Send:    lndSendClient,
			Receive: lndReceiveClient,
		},
	}

	ayo := aayoHttp.Service{
		Config:      config,
		DbClient:    dbClient,
		RedisClient: redisClient,
	}

	app := fiber.New()
	app.Use(logPkg.FiberRequestDebug)
	app.Use(networkPkg.CorsFiberMiddleware)
	app.Mount("/auth", auth.RegisterRoutes())
	app.Mount("/wallet", wallet.RegisterRoutes())
	app.Mount("/aayo", ayo.RegisterRoutes())

	err = app.Listen(fmt.Sprintf(":%v", config.AppPort))
	if err != nil {
		log.Fatalf("server listen failed, %s", err.Error())
	}
}
