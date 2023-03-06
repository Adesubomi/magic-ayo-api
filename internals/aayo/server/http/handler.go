package http

import (
	"fmt"
	aayoData "github.com/Adesubomi/magic-ayo-api/internals/aayo/data"
	authPkg "github.com/Adesubomi/magic-ayo-api/pkg/auth"
	configPkg "github.com/Adesubomi/magic-ayo-api/pkg/config"
	responsePkg "github.com/Adesubomi/magic-ayo-api/pkg/response"
	utilPkg "github.com/Adesubomi/magic-ayo-api/pkg/util"
	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Handler struct {
	Config      *configPkg.Config
	DbClient    *gorm.DB
	RedisClient *redis.Client
}

func (h Handler) StartGame(ctx *fiber.Ctx) error {
	userSession := authPkg.UserSessionFromFiberCtx(ctx)
	if userSession == nil {
		return responsePkg.Unauthorized(ctx, "user information not found")
	}

	newGame, err := h.getGamePlayRepo().
		StartGame(userSession.User.ID)

	if err != nil {
		resMessage := "Unable to start a new game -" + err.Error()
		return responsePkg.BadRequest(
			ctx,
			resMessage)
	}

	gameState, _ := newGame.GetLatestGameState()
	return responsePkg.Success(
		ctx,
		"New game started",
		map[string]interface{}{
			"gamePlay":  newGame,
			"gameState": gameState,
		})
}

func (h Handler) PotPack(ctx *fiber.Ctx) error {
	type r struct {
		PotIndex int `json:"pot"`
	}

	input := new(r)
	err := ctx.BodyParser(&input)
	if err != nil {

		return responsePkg.BadRequest(
			ctx,
			"check your inputs")
	}

	validationErrors := utilPkg.TranslateFiberValidationErrors(*input)
	if validationErrors != nil {
		return responsePkg.UnprocessableEntity(
			ctx,
			"Invalid Input",
			validationErrors,
		)
	}

	userSession := authPkg.UserSessionFromFiberCtx(ctx)
	if userSession == nil {
		return responsePkg.Unauthorized(
			ctx,
			"user information not found")
	}

	myGame, err := h.getGamePlayRepo().
		GetActiveGame(userSession.User.ID)
	if err != nil {
		return responsePkg.BadRequest(
			ctx, "unable to fetch active game")
	}

	err = h.getGamePlayRepo().MakeMove(myGame, input.PotIndex)
	if err != nil {
		return responsePkg.BadRequest(
			ctx,
			err.Error())
	}

	fmt.Println("<........>")
	fmt.Println(myGame.Moves)

	gameState, _ := myGame.GetLatestGameState()
	return responsePkg.Success(
		ctx,
		"Pot pack played",
		map[string]interface{}{
			"lastMove":  input.PotIndex,
			"gamePlay":  myGame,
			"gameState": gameState,
		})
}

func (h Handler) AbortGame(*fiber.Ctx) error {

	return nil
}

func (h Handler) GameStatus(ctx *fiber.Ctx) error {
	userSession := authPkg.UserSessionFromFiberCtx(ctx)
	if userSession == nil {
		return responsePkg.Unauthorized(
			ctx,
			"user information not found")
	}

	myGame, err := h.getGamePlayRepo().
		GetActiveGame(userSession.User.ID)
	if err != nil {
		return responsePkg.BadRequest(
			ctx, "unable to fetch active game")
	}

	gameState, _ := myGame.GetLatestGameState()
	return responsePkg.Success(
		ctx,
		"Game status - Aayo",
		map[string]interface{}{
			"gamePlay":  myGame,
			"gameState": gameState,
		})
}

func (h Handler) getGamePlayRepo() *aayoData.Repo {
	return &aayoData.Repo{
		Config:      h.Config,
		DbClient:    h.DbClient,
		RedisClient: h.RedisClient,
	}
}
