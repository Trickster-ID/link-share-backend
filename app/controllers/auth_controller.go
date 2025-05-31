package controllers

import (
	"github.com/gofiber/fiber/v3"
	"linkshare/app/dto"
	"linkshare/app/global/helper"
	"linkshare/app/global/model"
	"linkshare/app/usecases"
)

type IAuthController interface {
	Login(f fiber.Ctx) error
	RefreshToken(f fiber.Ctx) error
}

type authController struct {
	authUseCase usecases.IAuthUseCase
}

func NewAuthController(authUseCase usecases.IAuthUseCase) IAuthController {
	return &authController{
		authUseCase: authUseCase,
	}
}

func (c *authController) Login(f fiber.Ctx) error {
	response := &model.BaseResponse{}
	loginModel := new(dto.LoginRequest)
	if err := f.Bind().JSON(loginModel); err != nil {
		response.ErrorLog = helper.WriteLog(err, 400, "parsing body request is failed")
		return helper.Response(f, response)
	}
	response.Data, response.ErrorLog = c.authUseCase.Login(loginModel, f.Context())
	if response.ErrorLog != nil {
		return helper.Response(f, response)
	}
	return helper.Response(f, response)
}

func (c *authController) RefreshToken(f fiber.Ctx) error {
	response := &model.BaseResponse{}
	request := new(dto.RefreshTokenRequest)
	if err := f.Bind().JSON(request); err != nil {
		response.ErrorLog = helper.WriteLog(err, 400, "parsing body request is failed")
		return helper.Response(f, response)
	}
	response.Data, response.ErrorLog = c.authUseCase.RefreshToken(request, f.Context())
	if response.ErrorLog != nil {
		return helper.Response(f, response)
	}

	return helper.Response(f, response)
}
