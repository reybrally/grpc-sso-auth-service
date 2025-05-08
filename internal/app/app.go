package app

import (
	"github.com/reybrally/grpc-sso-auth-service/internal/app/grpc"
	"github.com/reybrally/grpc-sso-auth-service/internal/services/auth"
	"github.com/reybrally/grpc-sso-auth-service/internal/storage/sqlite"
	"log/slog"
	"time"
)

type App struct {
	GRPCsv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	TokenTTL time.Duration,
) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, TokenTTL)
	grpcApp := grpcapp.NewApp(log, authService, grpcPort)

	return &App{
		GRPCsv: grpcApp,
	}
}
