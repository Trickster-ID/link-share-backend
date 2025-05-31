package helper

import (
	"github.com/gofiber/fiber/v2"
	"linkshare/app/global/model"
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
