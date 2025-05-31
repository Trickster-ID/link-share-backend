package helper

import (
	"github.com/gofiber/fiber/v3"
	"linkshare/app/global/model"
	"net/http"
	"reflect"
)

func Response(ctx fiber.Ctx, response *model.BaseResponse, statusCodes ...int) error {
	statusCode := http.StatusOK
	if len(statusCodes) > 0 {
		statusCode = statusCodes[0]
	}

	if response == nil {
		return ctx.Status(statusCode).JSON(nil)
	}

	if response.ErrorLog != nil {
		response.TotalData = 0
		response.Data = nil
		response.StatusMessage = "error"
		return ctx.Status(response.ErrorLog.StatusCode).JSON(response)
	}

	if response.TotalData == 0 {
		totalData := int64(0)
		if response.Data != nil {
			v := reflect.ValueOf(response.Data)
			t := v.Type()

			// Check if the input is a slice or array
			if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
				totalData = int64(v.Len())
			} else {
				totalData = 1
			}
		}
		response.TotalData = totalData
	}
	response.StatusMessage = "success"
	return ctx.Status(statusCode).JSON(response)
}
