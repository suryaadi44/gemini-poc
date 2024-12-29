package logger

import (
	"gemini-poc/utils/config"

	"github.com/mattn/go-colorable"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"moul.io/zapfilter"
)

func InitLogger(conf *config.Config, namespace string) *zap.Logger {
	var level zapcore.Level
	switch conf.Log.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	}

	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   conf.Log.Path,
		MaxSize:    conf.Log.MaxSize,
		MaxBackups: conf.Log.MaxBackups,
		MaxAge:     conf.Log.MaxAge,
	})
	fileEncoderConfig := ecszap.NewDefaultEncoderConfig()
	fileCore := ecszap.NewCore(fileEncoderConfig, fileWriter, level)
	fileCore.With(
		[]zapcore.Field{zap.String("service", conf.App.Service), zap.String("environment", conf.App.Environment)})

	consoleWriter := zapcore.AddSync(colorable.NewColorableStdout())
	consoleEncoderConfig := zap.NewDevelopmentEncoderConfig()
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)
	consoleCore := zapcore.NewCore(consoleEncoder, consoleWriter, level)
	consoleRules := zapfilter.Reverse(zapfilter.MustParseRules("*:*access-log"))

	core := zapcore.NewTee(
		fileCore,
		zapfilter.NewFilteringCore(consoleCore, consoleRules),
	)

	logger := zap.New(core, zap.AddCaller())

	return logger.Named(namespace)
}
