package logger

import (
	"os"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/stretchr/testify/assert"
)

func TestNewLoggerConfig(t *testing.T) {
	tests := []struct {
		name  string
		level zap.AtomicLevel
		want  zapcore.Level
	}{
		{"Debug Level", zap.NewAtomicLevelAt(zap.DebugLevel), zap.DebugLevel},
		{"Info Level", zap.NewAtomicLevelAt(zap.InfoLevel), zap.InfoLevel},
		{"Warn Level", zap.NewAtomicLevelAt(zap.WarnLevel), zap.WarnLevel},
		{"Error Level", zap.NewAtomicLevelAt(zap.ErrorLevel), zap.ErrorLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := NewLoggerConfig(tt.level)
			assert.Equal(t, tt.want, cfg.Level.Level())
			assert.Equal(t, "json", cfg.Encoding)
			assert.Contains(t, cfg.OutputPaths, "stdout")
			assert.Contains(t, cfg.ErrorOutputPaths, "stderr")
		})
	}
}

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name      string
		envValue  string
		wantLevel zapcore.Level
	}{
		{"Default to Info", "", zap.InfoLevel},
		{"Set to Debug", "debug", zap.DebugLevel},
		{"Set to Info", "info", zap.InfoLevel},
		{"Set to Warn", "warn", zap.WarnLevel},
		{"Set to Error", "error", zap.ErrorLevel},
		{"Invalid Level", "invalid", zap.InfoLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable for the test
			os.Setenv("LOG_LEVEL", tt.envValue)
			defer os.Unsetenv("LOG_LEVEL")

			logger, err := NewLogger()
			assert.NoError(t, err)
			assert.NotNil(t, logger)
		})
	}
}
