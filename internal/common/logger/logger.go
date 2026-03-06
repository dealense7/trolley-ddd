package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dealense7/go-rates-ddd/internal/common/cfg"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ModuleName string

const (
	General ModuleName = "general"
)

type Factory interface {
	For(moduleName ModuleName) *zap.Logger
}

type factoryImpl struct {
	config    *cfg.Config
	baseLevel zapcore.Level
}

func NewFactory(cfg *cfg.Config) Factory {
	return &factoryImpl{
		config:    cfg,
		baseLevel: zap.InfoLevel,
	}
}

func (f *factoryImpl) For(moduleName ModuleName) *zap.Logger {

	moduleStr := string(moduleName)

	// Example: logs/2026-03-07/
	dateFolder := time.Now().Format("2006-01-02")
	logDir := filepath.Join("logs", dateFolder)

	// Ensure folder exists
	_ = os.MkdirAll(logDir, 0755)

	// Example: logs/2026-03-07/general.log
	fileName := filepath.Join(logDir, fmt.Sprintf("%s.log", moduleStr))

	fileRotator := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    10,
		MaxBackups: 7,
		MaxAge:     30,
		Compress:   true,
	}

	var writers []zapcore.WriteSyncer

	// Write to file unless testing
	if f.config.Env != "testing" {
		writers = append(writers, zapcore.AddSync(fileRotator))
	}

	// Write to console unless production
	if f.config.Env != "production" {
		writers = append(writers, zapcore.AddSync(os.Stdout))
	}

	coreWriter := zapcore.NewMultiWriteSyncer(writers...)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	encoder := zapcore.NewJSONEncoder(encoderConfig)

	core := zapcore.NewCore(encoder, coreWriter, f.baseLevel)

	return zap.New(core, zap.AddCaller()).With(zap.String("module", moduleStr))
}
