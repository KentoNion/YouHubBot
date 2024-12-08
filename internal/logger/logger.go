package logger

import (
	"YoutHubBot/internal/config"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func InitLogger(cfg *config.Config) (*zap.Logger, error) {
	// Открываем/создаём файл для логирования
	logFile := ""
	if cfg.ConfigPath != "" {
		logFile, err := os.OpenFile(cfg.ConfigPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open or create log file: %w", err)
		}
	}

	// Конфигурация в файл
	Encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	// Создаем core для вывода в консоль
	consoleCore := zapcore.NewCore(
		Encoder,
		zapcore.AddSync(os.Stdout), // Выводим в стандартный вывод (терминал)
		zapcore.DebugLevel,
	)

	// Вывод в файл
	fileCore := zapcore.NewCore(
		Encoder,
		zapcore.AddSync(logFile), // Выводим в лог-файл
		zapcore.DebugLevel,
	)

	// Хочу логи и в консоли и в файле
	core := zapcore.NewTee(consoleCore, fileCore)
	logger := zap.New(core)
	return logger, nil
}
