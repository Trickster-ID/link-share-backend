package controllers

import (
	"github.com/gofiber/fiber/v2"
	"linkshare/app/dto"
	"linkshare/app/global/helper"
	"linkshare/app/global/model"
	"linkshare/generated"
	"net/http"
	"strings"
)

func (c *Server) PostRegister(f *fiber.Ctx) error {
	response := &model.BaseResponse{}
	requestBody := &generated.RegisterRequest{}
	err := f.BodyParser(&requestBody)
	if err != nil {
		response.ErrorLog = helper.WriteLog(err, 400, "parsing body request is failed")
		return helper.Response(f, response)
	}
	if strings.TrimSpace(requestBody.Username) == "" || strings.TrimSpace(requestBody.Password) == "" || strings.TrimSpace(string(requestBody.Email)) == "" {
		response.ErrorLog = helper.WriteLog(err, 400, "username or password is empty")
		return helper.Response(f, response)
	}
	response.ErrorLog = c.authUseCase.Register(f.Context(), requestBody)
	if response.ErrorLog != nil {
		return helper.Response(f, response)
	}
	return helper.Response(f, response, http.StatusCreated)
}

func (c *Server) Login(f *fiber.Ctx) error {
	response := &model.BaseResponse{}
	loginModel := new(dto.LoginRequest)
	if err := f.BodyParser(&loginModel); err != nil {
		response.ErrorLog = helper.WriteLog(err, 400, "parsing body request is failed")
		return helper.Response(f, response)
	}
	response.Data, response.ErrorLog = c.authUseCase.Login(f.Context(), loginModel)
	if response.ErrorLog != nil {
		return helper.Response(f, response)
	}
	return helper.Response(f, response)
}

func (c *Server) RefreshToken(f *fiber.Ctx) error {
	response := &model.BaseResponse{}
	request := new(dto.RefreshTokenRequest)
	if err := f.BodyParser(&request); err != nil {
		response.ErrorLog = helper.WriteLog(err, 400, "parsing body request is failed")
		return helper.Response(f, response)
	}
	response.Data, response.ErrorLog = c.authUseCase.RefreshToken(f.Context(), request)
	if response.ErrorLog != nil {
		return helper.Response(f, response)
	}

	return helper.Response(f, response)
}
