package grpcapp

import (
	"fmt"
	authgrpc "github.com/reybrally/grpc-sso-auth-service/internal/grpc/auth"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	grpcServer *grpc.Server
	port       int
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}
func NewApp(
	log *slog.Logger,

	authService authgrpc.Auth,
	port int,
) *App {
	grpcServer := grpc.NewServer()
	authgrpc.Register(grpcServer, authService)
	return &App{
		log:        log,
		grpcServer: grpcServer,
		port:       port,
	}
}

func (a *App) Run() error {
	const op = "app.Start"
	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("grpc server running", slog.String("address", l.Addr().String()))
	if err := a.grpcServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *App) Stop() {
	const op = "app.Stop"
	a.log.With(slog.String("op", op)).Info("grpc server stopping", slog.Int("port", a.port))
	a.grpcServer.GracefulStop()
}
