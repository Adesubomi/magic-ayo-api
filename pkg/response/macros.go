package response

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// BadRequest : 400
func BadRequest(c *fiber.Ctx, msg string) error {
	if msg == "" {
		msg = "Bad Request"
	}

	return c.Status(http.StatusBadRequest).JSON(
		&fiber.Map{"message": msg},
	)
}

// Unauthorized : 401
func Unauthorized(c *fiber.Ctx, msg string) error {
	if msg == "" {
		msg = "You have not been signed in"
	}

	return c.Status(http.StatusUnauthorized).JSON(
		&fiber.Map{"message": msg},
	)
}

// Forbidden : 403
func Forbidden(c *fiber.Ctx, msg string) error {
	if msg == "" {
		msg = "You are not authorized to perform this action"
	}

	return c.Status(http.StatusForbidden).JSON(
		&fiber.Map{"message": msg},
	)
}

// NotFound : 404
func NotFound(c *fiber.Ctx, msg string) error {
	if msg == "" {
		msg = "Resource not found"
	}

	return c.Status(http.StatusNotFound).JSON(
		&fiber.Map{"message": msg},
	)
}

// UnprocessableEntity : 422
func UnprocessableEntity(c *fiber.Ctx, msg string, errs map[string]string) error {
	if msg == "" {
		msg = "Invalid input"
	}

	return c.Status(http.StatusUnprocessableEntity).JSON(
		&fiber.Map{
			"message": msg,
			"errors":  errs,
		},
	)
}

// InternalServerError : 500
func InternalServerError(c *fiber.Ctx, msg string, err error) error {
	if msg == "" {
		msg = "Internal Server Error"
	}

	return c.Status(http.StatusInternalServerError).JSON(
		&fiber.Map{"message": msg},
	)
}

// Success : 200
func Success(c *fiber.Ctx, msg string, data map[string]interface{}) error {
	return c.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": msg,
			"data":    data,
		},
	)
}

// PartialSuccess : 206
func PartialSuccess(c *fiber.Ctx, msg string, data map[string]interface{}) error {
	return c.Status(http.StatusPartialContent).JSON(
		&fiber.Map{
			"message": msg,
			"data":    data,
		},
	)
}

func PlainText(c *fiber.Ctx, msg string) error {
	return c.Status(http.StatusOK).Send([]byte(msg))
}
