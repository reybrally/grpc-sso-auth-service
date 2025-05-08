package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/reybrally/grpc-sso-auth-service/internal/domain/models"
	"github.com/reybrally/grpc-sso-auth-service/internal/lib/jwt"
	"github.com/reybrally/grpc-sso-auth-service/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

var (
	InvalidCredentials = errors.New("invalid credentials")
	InvalidAppId       = errors.New("invalid app id")
)

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
}

type UserProvider interface {
	GetUser(ctx context.Context, email string) (user models.User, err error)
	IsAdmin(ctx context.Context, userId int64) (isAdmin bool, err error)
}

type AppProvider interface {
	App(ctx context.Context, appId int) (app models.App, err error)
}

type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	TokenTTL    time.Duration
}

// New returns the instance of Auth service
func New(
	log *slog.Logger,
	usrSaver UserSaver,
	usrProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration) *Auth {
	return &Auth{
		log:         log,
		usrSaver:    usrSaver,
		usrProvider: usrProvider,
		appProvider: appProvider,
		TokenTTL:    tokenTTL,
	}
}

func (a *Auth) Login(ctx context.Context, email string, password string, appId int) (string, error) {
	const op = "Auth.Login"
	log := a.log.With(
		slog.String("op", op),
		//slog.String("email", email),
	)
	log.Info("attempting to login user")

	user, err := a.usrProvider.GetUser(ctx, email)
	if err != nil {
		if errors.Is(
			err,
			storage.ErrUserNotFound,
		) {
			a.log.Warn("user not found")
			return "", fmt.Errorf("%s: %w", op, InvalidCredentials)
		}
		a.log.Error("Failed to get user", err.Error())

		return "", fmt.Errorf("%s: %w", op, err)
	}
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", err.Error())
		return "", fmt.Errorf("%s: %w", op, InvalidCredentials)
	}
	app, err := a.appProvider.App(ctx, appId)

	if err != nil {
		a.log.Error("Failed to get app", err.Error())
		return "", fmt.Errorf("%s: %w", op, err)
	}
	log.Info("successfully logged in")

	token, err := jwt.NewToken(user, app, a.TokenTTL)
	if err != nil {
		a.log.Error("Failed to create token", err.Error())
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *Auth) RegisterNewUser(ctx context.Context, email string, password string) (int64, error) {
	const op = "register new user"
	log := a.log.With(slog.String("op", op)) // slog.String("email", email),

	log.Info("registering new user")
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash password", slog.String("error:", err.Error()))
		return 0, fmt.Errorf("%s : %w", op, err)
	}

	id, err := a.usrSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists")
			return 0, fmt.Errorf("%s: %w", op, err)
		}
		log.Error("failed to save user", slog.String("error:", err.Error()))
		return 0, fmt.Errorf("%s : %w", op, err)
	}

	log.Info("user created", slog.String("email", email))
	return id, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	const op = "IsAdmin"
	log := a.log.With(slog.String("op", op), slog.Int64("userId", userId))
	log.Info("checking if user is admin")
	isAdmin, err := a.usrProvider.IsAdmin(ctx, userId)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("user not found")
			return false, InvalidAppId
		}
		log.Error("failed to check if user is admin", slog.String("error:", err.Error()))
		return false, fmt.Errorf("%s : %w", op, err)
	}
	log.Info("checked if user is admin", slog.Bool("isAdmin", isAdmin))

	return isAdmin, nil
}
