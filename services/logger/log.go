package logger

import (
	"fmt"
	"github.com/beastars1/lol-prophet-gui/global"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Debug(msg string, keysAndValues ...interface{}) {
	global.Logger.Debugw(msg, keysAndValues...)
	//lol_prophet_gui.Append(msg)
}
func Info(msg string, keysAndValues ...interface{}) {
	global.Logger.Infow(msg, keysAndValues...)
	//lol_prophet_gui.Append(msg)
}
func Warn(msg string, keysAndValues ...interface{}) {
	//lol_prophet_gui.Append(msg)
	go sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetLevel(sentry.LevelWarning)
		scope.SetExtra("kv", keysAndValues)
		sentry.CaptureMessage(msg)
	})
	//global.Logger.Warnw(msg, keysAndValues...)
}
func Error(msg string, keysAndValues ...interface{}) {
	var errMsg string
	var errVerbose string
	for _, v := range keysAndValues {
		if field, ok := v.(zap.Field); ok && field.Type == zapcore.ErrorType {
			errMsg = field.Interface.(error).Error()
			errVerbose = fmt.Sprintf("%+v", field.Interface.(error))
		}
	}
	//lol_prophet_gui.Append(errVerbose)
	go sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetLevel(sentry.LevelError)
		scope.SetExtra("kv", keysAndValues)
		if errMsg != "" {
			scope.SetExtra("error", errMsg)
			scope.SetExtra("errorVerbose", errVerbose)
		}
		sentry.CaptureMessage(msg)
	})
	global.Logger.Errorw(msg, keysAndValues...)
}
