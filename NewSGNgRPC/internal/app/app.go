package app

import (
	grpcapp "NewSGNgRPC/internal/app/grpc"
	"log/slog"
)

type App struct {
	GRPCSrvr *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string) *App {
	// TODO: инициализация БД
	// TODO: инициализация сервиса
	grpcApp := grpcapp.New(log, grpcPort)
	return &App{
		GRPCSrvr: grpcApp,
	}
}
