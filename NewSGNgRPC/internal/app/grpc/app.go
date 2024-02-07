package grpcapp

import (
	FBgrpc "NewSGNgRPC/internal/grpc/feedback"
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, port int) *App {
	gRPCServer := grpc.NewServer()
	FBgrpc.Register(gRPCServer)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() error {
	if err := a.Run(); err != nil {
		panic(err)
	}
	return nil
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(slog.String("current operation", op), slog.Int("port", a.port))

	l, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("running gRPC-server...", slog.String("address", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *App) Stop() error {
	const op = "grpcapp.Run"

	log := a.log.With(slog.String("current operation", op), slog.Int("port", a.port))
	log.Info("stopping gRPC-server...")

	a.gRPCServer.GracefulStop()

	return nil
}
