package logging

import (
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	uzap "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

const (
	// For why the distinction between logr and zap levels, see: https://github.com/operator-framework/operator-sdk/issues/4771
	InfoLevel  = 0
	LogrDebugLevel = -1
	LogrTraceLevel = -2

	ZapDebugLevel = 1
	ZapTraceLevel = 2
)

// GetDevelopmentLogger returns a json logger configured for production. Level must be a negative number.
func GetDevelopmentLogger(level int, disableStacktrace bool) logr.Logger {
	if level > 0 {
		panic("logr logging levels cannot be positive")
	}
	zapLevel := zapcore.Level(level)
	atomicLevel := uzap.NewAtomicLevel()
	encoderCfg := uzap.NewDevelopmentEncoderConfig()

	encoderCfg.EncodeLevel = traceLevelFunc
	logger := uzap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atomicLevel,
	))
	atomicLevel.SetLevel(zapLevel)

	if !disableStacktrace {
		logger = logger.WithOptions(uzap.AddStacktrace(zapcore.Level(2)))
	}
	return zapr.NewLogger(logger)
}

// GetProductionLogger returns a console logger configured for development. Level must be a negative number.
func GetProductionLogger(level int, disableStacktrace bool) logr.Logger {
	if level > 0 {
		panic("logr logging levels cannot be positive")
	}
	// Set logging level
	zapLevel := zapcore.Level(level)
	atomicLevel := uzap.NewAtomicLevel()
	atomicLevel.SetLevel(zapLevel)
	// Create a production encoder
	encoderCfg := uzap.NewProductionEncoderConfig()
	encoderCfg.EncodeLevel = traceLevelFunc
	encoderCfg.EncodeTime = func(ts time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(ts.UTC().Format(time.RFC3339Nano))
	}
	// Crete a new json encoder, print to stdout at given atomic level
	zapCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atomicLevel,
	)
	// Configure sampling
	zapCore = zapcore.NewSamplerWithOptions(zapCore, time.Second, 100, 100)
	logger := uzap.New(zapCore)
	if !disableStacktrace {
		logger = logger.WithOptions(uzap.AddStacktrace(zapcore.Level(2)))
	}
	return zapr.NewLogger(logger)
}

// traceLevelFunc configures a zapcore.LevelEncoder to print "trace" as level field value in log outputs.
func traceLevelFunc(lvl zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
	if lvl == LogrTraceLevel {
		encoder.AppendString("trace")
	} else {
		encoder.AppendString(lvl.String())
	}
}