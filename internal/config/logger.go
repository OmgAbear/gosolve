package config

import (
	"log/slog"
	"os"
	"sync"
)

var (
	loggerInstance *slog.Logger
	loggerOnce     sync.Once
)

const (
	logLevelDebugStr = "Debug"
	logLevelInfoStr  = "Info"
	logLevelErrorStr = "Error"
)

func GetLogger() *slog.Logger {
	loggerOnce.Do(func() {
		config, _ := GetConfig()
		var logLevel slog.Level
		switch config.Logging.Level {
		case logLevelDebugStr:
			logLevel = slog.LevelDebug
		case logLevelInfoStr:
			logLevel = slog.LevelInfo
		case logLevelErrorStr:
			logLevel = slog.LevelError
		default:
			logLevel = slog.LevelInfo
		}

		loggerInstance = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		}))
		return
	})

	return loggerInstance
}
