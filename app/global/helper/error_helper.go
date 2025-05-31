package helper

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"linkshare/app/global/model"
	"net/http"
	"runtime"
	"strings"
)

var DefaultStatusText = map[int]string{
	http.StatusInternalServerError: "Terjadi Kesalahan, Silahkan Coba lagi Nanti",
	http.StatusNotFound:            "Data tidak Ditemukan",
	http.StatusBadRequest:          "Ada kesalahan pada request data, silahkan dicek kembali",
}

func WriteLog(err error, errorCode int, message interface{}) *model.ErrorLog {
	return writeLogRaw(err, errorCode, message, true)
}
func WriteLogWoP(err error, errorCode int, message interface{}) *model.ErrorLog {
	return writeLogRaw(err, errorCode, message, false)
}
func writeLogRaw(err error, errorCode int, message interface{}, isPrint bool) *model.ErrorLog {
	if pc, file, line, ok := runtime.Caller(2); ok {
		file = file[strings.LastIndex(file, "/")+1:]
		funcName := runtime.FuncForPC(pc).Name()
		output := &model.ErrorLog{
			StatusCode: errorCode,
			Err:        err,
		}
		outputForPrint := &model.ErrorLog{
			StatusCode:    errorCode,
			Err:           err,
			Line:          fmt.Sprintf("%d", line),
			Filename:      file,
			Function:      funcName,
			SystemMessage: err.Error(),
		}

		output.SystemMessage = err.Error()
		if message == nil {
			output.Message = DefaultStatusText[errorCode]
			if output.Message == "" {
				output.Message = http.StatusText(errorCode)
			}
		} else {
			output.Message = message
		}
		if errorCode == http.StatusInternalServerError {
			output.Line = fmt.Sprintf("%d", line)
			output.Filename = file
			output.Function = funcName
		}

		if isPrint {
			logForPrint := map[string]interface{}{}
			_ = DecodeMapType(outputForPrint, &logForPrint)
			logrus.SetReportCaller(false)
			logrus.WithFields(logForPrint).Error(message)
			logrus.SetReportCaller(true)
		}

		log := map[string]interface{}{}
		_ = DecodeMapType(output, &log)
		return output
	}

	return nil
}
