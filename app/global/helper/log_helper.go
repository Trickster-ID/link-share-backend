package helper

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

func GetCaller(skip int) string {
	// skip add +1 to ignore counting from under get caller
	if _, file, line, ok := runtime.Caller(skip + 1); ok {
		return fmt.Sprintf("%s:%d", file[strings.LastIndex(file, "/")+1:], line)
	}
	logrus.Error("Error getting caller")
	return ""
}
