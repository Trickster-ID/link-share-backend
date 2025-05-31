package controllers

import (
	"github.com/gofiber/fiber/v2"
	"linkshare/app/global/helper"
	"linkshare/app/global/model"
	"linkshare/generated"
)

// ensure that we've conformed to the `ServerInterface` with a compile-time check
var _ generated.ServerInterface = (*Server)(nil)

type Server struct{}

func NewServer() Server {
	return Server{}
}

// (GET /ping)
func (Server) Ping(ctx *fiber.Ctx) error {
	response := &model.BaseResponse{}
	response.Data = "pong"
	return helper.Response(ctx, response)
}
