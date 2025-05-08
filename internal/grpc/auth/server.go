package auth

import (
	"errors"
	"github.com/reybrally/grpc-sso-auth-service/internal/services/auth"
	"github.com/reybrally/grpc-sso-auth-service/internal/storage"
	sso "github.com/reybrally/protos/gen/go/sso"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context,
		email string,
		password string,
		userID int,
	) (token string, err error)

	RegisterNewUser(ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)

	IsAdmin(ctx context.Context,
		UserID int64,
	) (bool, error)
}

type serverAPI struct {
	sso.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	sso.RegisterAuthServer(
		gRPC,
		&serverAPI{auth: auth},
	)
}

const (
	emptyValue = 0
)

func (s *serverAPI) Login(ctx context.Context, req *sso.LoginRequest) (*sso.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetUserId()))
	if err != nil {
		if errors.Is(err, auth.InvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "internal error")

	}

	return &sso.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *sso.RegisterRequest) (*sso.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	UserID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		// TODO: ...
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "internal error")

	}
	return &sso.RegisterResponse{
		UserId: UserID,
	}, nil
}

func (s *serverAPI) IsAdminRequest(ctx context.Context, req *sso.IsAdminRequest) (*sso.IsAdminResponse, error) {
	if err := validateIsAdmin(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &sso.IsAdminResponse{IsAdmin: isAdmin}, nil
}

func validateLogin(req *sso.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "login required")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "login required")
	}
	if req.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "login_id is required")
	}
	return nil

}

func validateRegister(req *sso.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "login required")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "login required")
	}
	return nil
}

func validateIsAdmin(req *sso.IsAdminRequest) error {
	if req.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "login_id is required")
	}
	return nil
}
