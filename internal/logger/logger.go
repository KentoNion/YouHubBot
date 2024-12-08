package logger

import (
	"YoutHubBot/internal/config"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func InitLogger(cfg *config.Config) (*zap.Logger, error) {
	// Конфигурация в файл
	Encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	// Открываем/создаём файл для логирования
	// Создаем core для вывода в консоль
	consoleCore := zapcore.NewCore(
		Encoder,
		zapcore.AddSync(os.Stdout), // Выводим в стандартный вывод (терминал)
		zapcore.DebugLevel,
	)
	core := zapcore.NewTee(consoleCore)

	if cfg.ConfigPath != "" { //если передан параметр пути до файлов с логами, то добавляем в core этот файл
		logFile, err := os.OpenFile(cfg.ConfigPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open or create log file: %w", err)
		}
		// Вывод в файл
		fileCore := zapcore.NewCore(
			Encoder,
			zapcore.AddSync(logFile), // Выводим в лог-файл
			zapcore.DebugLevel,
		)
		core = zapcore.NewTee(consoleCore, fileCore)
	}

	logger := zap.New(core)
	return logger, nil
}
