package helper

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"linkshare/app/constants"
	"linkshare/app/global/model"
	"linkshare/app/models"
	"net/http"
	"reflect"
)

func Response(ctx *fiber.Ctx, response *model.BaseResponse, statusCodes ...int) error {
	statusCode := http.StatusOK
	if len(statusCodes) > 0 {
		statusCode = statusCodes[0]
	}

	// nil response case
	if response == nil {
		return ctx.Status(statusCode).JSON(nil)
	}

	// error case
	if err := response.ErrorLog; err != nil {
		response.StatusMessage = "error"
		response.TotalData = 0
		response.Data = nil
		return ctx.Status(err.StatusCode).JSON(response)
	}

	// only reflect if Data is not nil and TotalData is 0
	if response.TotalData == 0 && response.Data != nil {
		switch v := response.Data.(type) {
		case []any:
			response.TotalData = int64(len(v))
		case []string:
			response.TotalData = int64(len(v))
		case []int:
			response.TotalData = int64(len(v))
		default:
			// fallback to reflection only if not basic slice
			val := reflect.ValueOf(response.Data)
			if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
				response.TotalData = int64(val.Len())
			}
		}
	}

	response.StatusMessage = "success"
	return ctx.Status(statusCode).JSON(response)
}

func SetUserDataOnCtx(f *fiber.Ctx, data *models.UserDataOnJWT) error {
	f.Locals(constants.UserDataJwt, data)
	return nil
}

func GetUserDataOnCtx(ctx any) (*models.UserDataOnJWT, error) {
	var result *models.UserDataOnJWT
	if fctx, ok := ctx.(*fiber.Ctx); ok {
		result = fctx.Locals(constants.UserDataJwt).(*models.UserDataOnJWT)
	} else if newCtx, ok := ctx.(context.Context); ok {
		result = newCtx.Value(constants.UserDataJwt).(*models.UserDataOnJWT)
	} else {
		return nil, errors.New("invalid context type")
	}
	if result == nil {
		return nil, errors.New("user data not found")
	}
	return result, nil
}
