package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"linkshare/generated"
	"net/http"
)

// ensure that we've conformed to the `ServerInterface` with a compile-time check
var _ generated.ServerInterface = (*Server)(nil)

type Server struct{}

func NewServer() Server {
	return Server{}
}

// (GET /ping)
func (Server) GetPing(ctx *fiber.Ctx) error {
	resp := generated.Pong{
		Ping: "pong",
	}
	fmt.Println("hitted!")
	return ctx.Status(http.StatusOK).JSON(resp)
}
