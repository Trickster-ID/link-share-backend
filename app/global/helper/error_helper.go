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
	http.StatusInternalServerError: "something went wrong, try again later",
	http.StatusNotFound:            "data not found",
	http.StatusBadRequest:          "there something wrong with your request, please check your request",
}

func WriteLog(err error, errorCode int, message string) *model.ErrorLog {
	return writeLogRaw(err, errorCode, message, true)
}
func WriteLogWoP(err error, errorCode int, message string) *model.ErrorLog {
	return writeLogRaw(err, errorCode, message, false)
}
func writeLogRaw(err error, errorCode int, message string, isPrint bool) *model.ErrorLog {
	if pc, file, line, ok := runtime.Caller(2); ok {
		var errString string
		if err != nil {
			errString = err.Error()
		}
		file = file[strings.LastIndex(file, "/")+1:]
		funcName := runtime.FuncForPC(pc).Name()
		output := &model.ErrorLog{
			StatusCode: errorCode,
			Err:        err,
		}

		output.SystemMessage = errString
		output.Message = message
		if errorCode == http.StatusInternalServerError {
			output.Line = fmt.Sprintf("%d", line)
			output.Filename = file
			output.Function = funcName
			output.Message = fmt.Sprintf("%s", DefaultStatusText[errorCode])
		} else if errorCode == http.StatusNotFound || errorCode == http.StatusBadRequest {
			output.Message = fmt.Sprintf("%s", DefaultStatusText[errorCode])
		}

		if isPrint {
			outputForPrint := &model.ErrorLog{
				StatusCode:    errorCode,
				Err:           err,
				Line:          fmt.Sprintf("%d", line),
				Filename:      file,
				Function:      funcName,
				SystemMessage: errString,
			}
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
