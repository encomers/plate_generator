package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New возвращает настроенный *zap.Logger
func New(opts ...Option) *zap.Logger {
	cfg := &config{
		level:      "info",
		production: true,
		withCaller: true,
		withStack:  true,
		appName:    "my-app",
		version:    "1.0.0",
	}

	for _, opt := range opts {
		opt(cfg)
	}

	// Базовая конфигурация
	var zapCfg zap.Config
	if cfg.production {
		zapCfg = zap.NewProductionConfig()
	} else {
		zapCfg = zap.NewDevelopmentConfig()
		zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Красивые настройки
	encCfg := zap.NewProductionEncoderConfig()
	encCfg.TimeKey = "timestamp"
	encCfg.EncodeTime = zapcore.ISO8601TimeEncoder // 2026-05-13T22:45:12.345Z
	encCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	encCfg.EncodeDuration = zapcore.StringDurationEncoder
	encCfg.EncodeCaller = zapcore.ShortCallerEncoder

	zapCfg.EncoderConfig = encCfg
	zapCfg.Level = zap.NewAtomicLevelAt(parseLevel(cfg.level))
	zapCfg.OutputPaths = []string{"stdout"}
	zapCfg.ErrorOutputPaths = []string{"stderr"}

	// Начальные поля
	if len(cfg.initialFields) == 0 {
		zapCfg.InitialFields = map[string]any{
			"app":     cfg.appName,
			"version": cfg.version,
			"pid":     os.Getpid(),
		}
	} else {
		zapCfg.InitialFields = cfg.initialFields
	}

	logger, err := zapCfg.Build(
		zap.AddCallerSkip(1),
	)
	if err != nil {
		panic("failed to create logger: " + err.Error())
	}

	// Дополнительные опции
	if cfg.withCaller {
		logger = logger.WithOptions(zap.AddCaller())
	}
	if cfg.withStack {
		logger = logger.WithOptions(zap.AddStacktrace(zapcore.ErrorLevel))
	}

	return logger
}

// parseLevel преобразует строку в zap уровень
func parseLevel(level string) zapcore.Level {
	var l zapcore.Level
	if err := l.UnmarshalText([]byte(level)); err != nil {
		return zapcore.InfoLevel
	}
	return l
}

// ======================= Опции =======================

type config struct {
	level         string
	production    bool
	withCaller    bool
	withStack     bool
	appName       string
	version       string
	initialFields map[string]any
}

type Option func(*config)

func WithLevel(level string) Option {
	return func(c *config) { c.level = level }
}

func WithDevelopment() Option {
	return func(c *config) { c.production = false }
}

func WithProduction() Option {
	return func(c *config) { c.production = true }
}

func WithCaller(enable bool) Option {
	return func(c *config) { c.withCaller = enable }
}

func WithStacktrace(enable bool) Option {
	return func(c *config) { c.withStack = enable }
}

func WithAppName(name string) Option {
	return func(c *config) { c.appName = name }
}

func WithVersion(version string) Option {
	return func(c *config) { c.version = version }
}

func WithInitialFields(fields map[string]any) Option {
	return func(c *config) { c.initialFields = fields }
}
