package http

import (
	"errors"
	"fmt"
	userData "github.com/Adesubomi/magic-ayo-api/internals/user/data"
	authPkg "github.com/Adesubomi/magic-ayo-api/pkg/auth"
	configPkg "github.com/Adesubomi/magic-ayo-api/pkg/config"
	logPkg "github.com/Adesubomi/magic-ayo-api/pkg/log"
	responsePkg "github.com/Adesubomi/magic-ayo-api/pkg/response"
	utilPkg "github.com/Adesubomi/magic-ayo-api/pkg/util"
	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strings"
)

type Handler struct {
	Config      *configPkg.Config
	DbClient    *gorm.DB
	RedisClient *redis.Client
}

func (h Handler) getUserRepo() userData.UserRepo {
	return userData.UserRepo{
		Config:      h.Config,
		DbClient:    h.DbClient,
		RedisClient: h.RedisClient,
	}
}

func (h Handler) getAuthRegistry() *authPkg.Registry {
	return &authPkg.Registry{
		RedisClient: h.RedisClient,
	}
}

func (h Handler) SignUp(ctx *fiber.Ctx) error {
	type signupRequest struct {
		ID       string `json:"id"`
		Password string `json:"password"`
	}

	signupInput := new(signupRequest)
	err := ctx.BodyParser(&signupInput)
	if err != nil {

		return responsePkg.BadRequest(
			ctx,
			"Unable to process your request, check your inputs")
	}

	validationErrors := utilPkg.TranslateFiberValidationErrors(*signupInput)
	if validationErrors != nil {
		return responsePkg.UnprocessableEntity(
			ctx,
			"Invalid Input",
			validationErrors,
		)
	}

	user, err := h.getUserRepo().CreateAccount(
		signupInput.ID,
		signupInput.Password)
	if err != nil && errors.Is(err, logPkg.DuplicateRecordError) {
		return responsePkg.UnprocessableEntity(
			ctx,
			"Sign up failed.",
			map[string]string{
				"id": "User ID has already been used",
			})
	} else if err != nil {
		logPkg.ReportError(err)
		return responsePkg.BadRequest(
			ctx,
			"Unexpected error occurred. Sign Up failed",
		)
	}

	// generate token and claim
	userSession, err := h.getAuthRegistry().
		CreateUserSession(user)
	if err != nil {
		return responsePkg.InternalServerError(ctx,
			"unable to log you in at the moment",
			err)
	}

	return responsePkg.Success(
		ctx,
		"Sign up successful.",
		map[string]interface{}{
			"user":    user,
			"session": userSession.SerializeForResponse(),
		})
}

func (h Handler) Login(ctx *fiber.Ctx) error {
	type loginRequest struct {
		ID       string `json:"id"`
		Password string `json:"password"`
	}

	input := new(loginRequest)
	err := ctx.BodyParser(&input)
	if err != nil {

		return responsePkg.BadRequest(
			ctx,
			"Unable to process your request, check your input")
	}

	validationErrors := utilPkg.TranslateFiberValidationErrors(*input)
	if validationErrors != nil {
		return responsePkg.UnprocessableEntity(
			ctx,
			"Invalid Input",
			validationErrors,
		)
	}

	credentialsErr := responsePkg.UnprocessableEntity(
		ctx,
		"Invalid data supplied",
		map[string]string{
			"id": fmt.Sprintf("Invalid ID or password"),
		},
	)
	identifier := strings.ToLower(input.ID)
	user, err := h.getUserRepo().FindUser(identifier)
	if err != nil && errors.Is(err, logPkg.RecordNotFoundError) {
		return credentialsErr
	} else if err != nil {
		fmt.Println(" - - error:", err.Error())
		return credentialsErr

	}

	// match password hash
	if !utilPkg.BcryptCompare(user.Password, input.Password) {
		return credentialsErr
	}

	// generate token and claim
	userSession, err := h.getAuthRegistry().
		CreateUserSession(user)
	if err != nil {
		return responsePkg.InternalServerError(ctx,
			"unable to log you in at the moment",
			err)
	}

	return responsePkg.Success(
		ctx,
		"Login successful",
		map[string]interface{}{
			"user":    user,
			"session": userSession.SerializeForResponse(),
		})
}
