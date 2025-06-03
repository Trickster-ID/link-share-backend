package controllers

import (
	"github.com/gofiber/fiber/v2"
	"linkshare/app/global/helper"
	"linkshare/app/global/model"
	"linkshare/app/usecases"
	"linkshare/generated"
)

// ensure that we've conformed to the `ServerInterface` with a compile-time check
var _ generated.ServerInterface = (*Server)(nil)

type Server struct {
	authUseCase usecases.IAuthUseCase
}

func NewServer(authUseCase usecases.IAuthUseCase) *Server {
	return &Server{
		authUseCase: authUseCase,
	}
}

// (GET /ping)
func (*Server) Ping(ctx *fiber.Ctx) error {
	response := &model.BaseResponse{}
	response.Data = "pong"
	return helper.Response(ctx, response)
}
