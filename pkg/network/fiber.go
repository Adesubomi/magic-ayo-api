package network

import "github.com/gofiber/fiber/v2/middleware/cors"

var CorsFiberMiddleware = cors.New(
	cors.Config{
		Next:             nil,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	})
