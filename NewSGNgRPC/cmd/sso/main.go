package main

import (
	"NewSGNgRPC/internal/app"
	"NewSGNgRPC/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	//инициализация объекта конфига
	cfg := config.MustLoad()
	//инициализация логгера
	log := setupLogger(envLocal)
	log.Info("start of app - env", slog.String("env", cfg.Env))
	log.Info("start of app - db_path", slog.String("db_path", cfg.StoragePath))
	log.Info("start of app - port", slog.Int("port", cfg.GRPC.Port))

	// TODO: блядь просто посмотреть видео и не заниматься хуйней
	//FBService = FBgrpc.Feedback()

	// инициализация приложения (package app)
	// TODO: token TTL - не нужен в проекте, грамотно почистить
	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath)
	go application.GRPCSrvr.MustRun()

	//Graceful Shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop
	application.GRPCSrvr.Stop()
	log.Info("application stopped.", slog.Any("Signal", sign))
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
