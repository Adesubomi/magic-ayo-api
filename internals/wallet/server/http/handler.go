package http

import (
	"context"
	"fmt"
	"github.com/Adesubomi/magic-ayo-api/internals/wallet/data"
	"github.com/Adesubomi/magic-ayo-api/internals/wallet/entity"
	authPkg "github.com/Adesubomi/magic-ayo-api/pkg/auth"
	configPkg "github.com/Adesubomi/magic-ayo-api/pkg/config"
	lightningPkg "github.com/Adesubomi/magic-ayo-api/pkg/lightning"
	responsePkg "github.com/Adesubomi/magic-ayo-api/pkg/response"
	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/lightningnetwork/lnd/lnrpc"
	"gorm.io/gorm"
	"time"
)

type Handler struct {
	Config      *configPkg.Config
	DbClient    *gorm.DB
	RedisClient *redis.Client
	LNClients   *lightningPkg.LNClients
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

func (h Handler) GenerateInvoice(ctx *fiber.Ctx) error {
	userSession := authPkg.UserSessionFromFiberCtx(ctx)
	amount := int64(20)
	amountMsat := amount * 1000
	// Create a new invoice

	//Generate the payment request URL
	response, err := h.LNClients.Receive.Client.AddInvoice(
		context.Background(),
		&lnrpc.Invoice{
			Memo:      "Ready Gamer ",
			ValueMsat: amountMsat,
		})
	if err != nil {
		fmt.Printf("Failed to generate payment request URL: %v\n", err)
		return responsePkg.BadRequest(ctx, "unable to generate invoice")
	}

	invoice := entity.LNInvoice{
		UserID:         userSession.User.ID,
		RequestHash:    string(response.RHash),
		PaymentRequest: response.GetPaymentRequest(),
		AddIndex:       response.GetAddIndex(),
	}

	// create an invoice in the db

	// Print the payment request URL
	return responsePkg.Success(
		ctx,
		"Lightning Invoice",
		map[string]interface{}{
			"invoice": map[string]string{
				"paymentRequest": invoice.PaymentRequest,
				"amount":         fmt.Sprintf("%v", amount),
			},
		})
}

func (h Handler) GetInvoiceStatus(ctx *fiber.Ctx) error {
	lnUrl := ctx.Params("url")
	paymentRequest, err := h.LNClients.Send.Client.DecodePayReq(
		context.Background(), &lnrpc.PayReqString{
			PayReq: lnUrl,
		})
	if err != nil {
		return responsePkg.BadRequest(
			ctx,
			"Invalid payment link")
	}

	fmt.Println(" - - - - - error::", paymentRequest.PaymentHash)
	invoice, err := h.LNClients.Send.Client.
		LookupInvoice(
			context.Background(),
			&lnrpc.PaymentHash{
				RHash: []byte(paymentRequest.PaymentHash),
			})
	if err != nil {
		fmt.Println(" - - - - - error::", err.Error())
		return responsePkg.BadRequest(
			ctx,
			"Invoice not found!")
	}

	return responsePkg.Success(
		ctx,
		"Invoice Info",
		map[string]interface{}{
			"invoice": map[string]string{
				"paymentRequest": invoice.PaymentRequest,
				"amount":         fmt.Sprintf("%v", invoice.ValueMsat),
			},
		})
}

func (h Handler) MakePayment(ctx *fiber.Ctx) error {
	lnUrl := ctx.Params("ln-url")
	paymentRequest, err := h.LNClients.Send.Client.DecodePayReq(
		context.Background(), &lnrpc.PayReqString{
			PayReq: lnUrl,
		})
	if err != nil {
		return responsePkg.BadRequest(
			ctx,
			"Invalid payment link")
	}

	invoice, err := h.LNClients.Send.Client.
		LookupInvoice(
			context.Background(),
			&lnrpc.PaymentHash{
				RHash: []byte(paymentRequest.PaymentHash),
			})
	if err != nil {
		return responsePkg.BadRequest(
			ctx,
			"Invoice not found!")
	}

	expiryTime := time.Unix(invoice.CreationDate+invoice.Expiry, 0)
	if time.Now().After(expiryTime) {
		return responsePkg.BadRequest(
			ctx,
			"Invoice has expired.")

	}

	payment, err := h.LNClients.Send.Client.SendPaymentSync(
		context.Background(),
		&lnrpc.SendRequest{
			PaymentRequest: invoice.PaymentRequest,
		})
	if err != nil {
		return responsePkg.BadRequest(
			ctx,
			"Payment failed")
	} else if payment.PaymentError != "" {
		fmt.Println("     [✗] payment error:", payment.PaymentError)
		return responsePkg.BadRequest(
			ctx,
			"Error occurred with payment")
	}

	fmt.Println("  think we've made payment", payment)

	return responsePkg.Success(
		ctx,
		"Payment made",
		map[string]interface{}{
			"payment": payment,
		},
	)
}
