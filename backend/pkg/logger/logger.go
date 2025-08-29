package logger

import (
	"os"
	"path/filepath"

	"github.com/Luawig/neoneuro/backend/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var base *zap.Logger

func parseLevel(lvl string) zapcore.Level {
	switch lvl {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func ensureDir(path string) {
	dir := filepath.Dir(path)
	_ = os.MkdirAll(dir, 0o755)
}

// Init creates a global structured logger using zap + lumberjack rotation.
// No prod/dev branch; single consistent JSON output to both file and stdout.
func Init(cfg config.Config) {
	if base != nil {
		return
	}

	level := parseLevel(cfg.Log.Level)
	encCfg := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// file sink with rotation
	ensureDir(cfg.Log.File)
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.Log.File,
		MaxSize:    cfg.Log.MaxSize,
		MaxBackups: cfg.Log.MaxBackups,
		MaxAge:     cfg.Log.MaxAge,
		Compress:   cfg.Log.Compress,
	})

	jsonEnc := zapcore.NewJSONEncoder(encCfg)

	core := zapcore.NewTee(
		zapcore.NewCore(jsonEnc, fileWriter, level),                 // file
		zapcore.NewCore(jsonEnc, zapcore.AddSync(os.Stdout), level), // stdout
	)

	base = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

func L() *zap.Logger {
	if base == nil {
		panic("logger not initialized")
	}
	return base
}

func S() *zap.SugaredLogger { return L().Sugar() }
