package http

import (
	"context"
	"fmt"
	"github.com/Adesubomi/magic-ayo-api/internals/wallet/data"
	authPkg "github.com/Adesubomi/magic-ayo-api/pkg/auth"
	configPkg "github.com/Adesubomi/magic-ayo-api/pkg/config"
	lightningPkg "github.com/Adesubomi/magic-ayo-api/pkg/lightning"
	responsePkg "github.com/Adesubomi/magic-ayo-api/pkg/response"
	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/lightningnetwork/lnd/lnrpc"
	"gorm.io/gorm"
	"log"
)

type Handler struct {
	Config      *configPkg.Config
	DbClient    *gorm.DB
	RedisClient *redis.Client
	LndClient   *lightningPkg.LNDClient
}

func (h Handler) getWalletRepo() *data.Repo {
	return &data.Repo{
		Config:      h.Config,
		DbClient:    h.DbClient,
		RedisClient: h.RedisClient,
	}
}

func (h Handler) GetUserWallet(ctx *fiber.Ctx) error {
	userSession := authPkg.UserSessionFromFiberCtx(ctx)
	if userSession == nil {
		return responsePkg.Unauthorized(ctx, "user information not found")
	}

	// get user wallet
	wallet, err := h.getWalletRepo().
		GetUserWallet(userSession.User.ID)

	if err != nil {
		return responsePkg.BadRequest(
			ctx,
			"unable to fetch user wallet")
	}

	return responsePkg.Success(
		ctx,
		"user wallet",
		map[string]interface{}{
			"wallet": wallet,
		},
	)
}

func (h Handler) GenerateInvoiceLNUrl(ctx *fiber.Ctx) error {
	amount := int64(2000)
	// Create a new invoice
	invoice := &lnrpc.Invoice{
		Memo:  "Example Memo",
		Value: amount,
	}

	// Generate the payment request URL
	response, err := h.LndClient.Client.AddInvoice(
		context.Background(),
		invoice)
	if err != nil {
		log.Fatalf("Failed to generate payment request URL: %v", err)
	}

	paymentRequest := response.String()

	// Print the payment request URL
	fmt.Printf("Payment Request: %s\n", paymentRequest)
	return responsePkg.Success(
		ctx,
		"Lightning invoice URL",
		map[string]interface{}{
			"lnUrl": paymentRequest,
		})
}

func (h Handler) GetLNInvoiceStatus(ctx *fiber.Ctx) error {
	return nil
}
