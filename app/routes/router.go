package routes

import (
	"github.com/gofiber/fiber/v2"
	"linkshare/app/controllers"
	"linkshare/app/global/helper"
	"linkshare/app/middleware"
)

func NewRouter(authController controllers.IAuthController) *fiber.App {
	f := fiber.New()

	f.Get("/ping", func(c *fiber.Ctx) error { return helper.Response(c, nil) })

	// region auth
	f.Post("/auth/login", authController.Login, middleware.BasicAuthMiddleware())
	f.Get("/auth/verify-token", func(c *fiber.Ctx) error { return helper.Response(c, nil) }, middleware.TokenMiddleware())
	f.Post("/auth/refresh-token", authController.RefreshToken, middleware.BasicAuthMiddleware())
	// endregion auth

	// region cms

	// endregion cms

	// region main page

	// endregion main page
	return f
}
