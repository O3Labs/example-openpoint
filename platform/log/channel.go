package log

import (
	"fmt"

	"github.com/o3labs/openpoint/platform/config"
	logrus "github.com/sirupsen/logrus"
	// "os"
	"time"

	"encoding/json"
	"net/http"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
)

type ChannelLogger struct {
	ID      uint `json:"-"`
	Message string
	Error   error
	Method  string
}

func Init(path string) {

	logrus.SetFormatter(&ChannelTextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true})

	if config.Env.Name != config.LocalEnv {

		writer, err := rotatelogs.New(
			path+".%Y-%m-%d",
			rotatelogs.WithLinkName(path),
			rotatelogs.WithMaxAge(24*time.Hour),
			rotatelogs.WithRotationTime(24*time.Hour),
		)

		if err != nil {
			Errorf("Failed to log to file because %+v", err)
			return
		}

		Info("Started log to file")

		logrus.SetOutput(writer)
		logrus.SetLevel(logrus.WarnLevel)
		logrus.SetFormatter(&ChannelJSONFormatter{})

	}

}

func Info(format string, args ...interface{}) {
	logrus.WithFields(logrus.Fields{}).Info(fmt.Sprintf(format, args...))
}

func Warn(format string, args ...interface{}) {
	logrus.WithFields(logrus.Fields{}).Warn(fmt.Sprintf(format, args...))
}

func Error(err error) {
	logrus.WithFields(logrus.Fields{"error": err}).Error()
}

func Errorf(format string, args ...interface{}) {
	logrus.WithFields(logrus.Fields{"error": fmt.Sprintf(format, args...)}).Error()
}

func ErrorHttpRequest(status int, duration time.Duration, request *http.Request) {
	// https://golang.org/src/net/http/request.go
	headerErr := ""
	header, err := json.Marshal(request.Header)
	if err == nil {
		headerErr = string(header)
	}
	logrus.WithFields(logrus.Fields{"error": fmt.Sprintf("%v", status), "duration": fmt.Sprintf("%v", duration), "method": request.Method, "requestURI": request.RequestURI, "header": headerErr}).Error()
}

func DebugMessage(format string, args ...interface{}) {
	logrus.WithFields(logrus.Fields{"dubug": fmt.Sprintf(format, args...)}).Debug()
}

func Panicf(format string, args ...interface{}) {
	logrus.WithFields(logrus.Fields{"error": fmt.Sprintf(format, args...)}).Error()
}

func Panic(err error) {
	// logrus panic Exit(1)
	logrus.WithFields(logrus.Fields{"panic": err}).Error()
}

func Fatal(err error) {
	// logrus fatal Exit(1)
	logrus.WithFields(logrus.Fields{"fatal": err}).Error()
}

func Printf(format string, args ...interface{}) {
	Info(fmt.Sprintf(format, args...))
}

func Println(format string, args ...interface{}) {
	Info(fmt.Sprintf(format, args...))
}
