package controllers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"linkshare/app/dto"
	"linkshare/app/global/helper"
	"linkshare/app/global/model"
	"linkshare/generated"
	"net/http"
	"strings"
)

func (c *Server) PostAuthRegister(f *fiber.Ctx) error {
	response := &model.BaseResponse{}
	requestBody := &generated.RegisterRequest{}
	err := f.BodyParser(&requestBody)
	if err != nil {
		response.ErrorLog = helper.WriteLog(err, 400, "parsing body request is failed")
		return helper.Response(f, response)
	}
	if strings.TrimSpace(requestBody.Username) == "" || strings.TrimSpace(requestBody.Password) == "" || strings.TrimSpace(string(requestBody.Email)) == "" {
		response.ErrorLog = helper.WriteLog(errors.New("username or password is empty"), 400, "")
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
	tokenString := strings.TrimSpace(strings.TrimPrefix(f.Get(fiber.HeaderAuthorization), "Bearer "))
	response.Data, response.ErrorLog = c.authUseCase.RefreshToken(f.Context(), &dto.RefreshTokenRequest{RefreshToken: tokenString})
	if response.ErrorLog != nil {
		return helper.Response(f, response)
	}

	return helper.Response(f, response)
}
