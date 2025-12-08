package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Log is the global logger instance
var Log *zap.SugaredLogger

// InitLogger initializes the global logger
// mode can be "development" or "production"
func InitLogger(mode string) {
	var config zap.Config

	if mode == "production" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Create file writer with rotation using lumberjack
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./logs/app.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     7, // days
		Compress:   true,
	})

	// Create console writer
	consoleWriter := zapcore.AddSync(os.Stdout)

	// Create encoder
	var encoder zapcore.Encoder
	if mode == "production" {
		encoder = zapcore.NewJSONEncoder(config.EncoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(config.EncoderConfig)
	}

	// Create core that writes to both console and file
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, consoleWriter, config.Level),
		zapcore.NewCore(
			zapcore.NewJSONEncoder(config.EncoderConfig),
			fileWriter,
			config.Level,
		),
	)

	// Create logger
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	// Set global logger as sugared logger for easier use
	Log = logger.Sugar()
}

// Sync flushes any buffered log entries
func Sync() {
	if Log != nil {
		Log.Sync()
	}
}
