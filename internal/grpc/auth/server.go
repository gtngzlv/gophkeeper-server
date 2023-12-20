package auth

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/gtngzlv/gophkeeper-protos/gen/go/gophkeeper"

	customerr "github.com/gtngzlv/gophkeeper-server/internal/domain/errors"
)

type IAuthService interface {
	Register(ctx context.Context, email string, password string) (userID int64, err error)
	Login(ctx context.Context, email string, password string) (token string, err error)
	Logout(ctx context.Context)
}

type serverAPI struct {
	pb.UnimplementedAuthServer

	authService IAuthService
}

func Register(grpcServer *grpc.Server, auth IAuthService) {
	pb.RegisterAuthServer(grpcServer, &serverAPI{
		authService: auth,
	})
}

func (s *serverAPI) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if err := validateRegister(in); err != nil {
		return nil, err
	}

	userID, err := s.authService.Register(ctx, in.GetEmail(), in.GetPassword())
	if err != nil {
		if errors.Is(err, customerr.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user with this email already exist")
		}
		return nil, status.Error(codes.Internal, "failed to register")
	}

	return &pb.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *serverAPI) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	if err := validateLogin(in); err != nil {
		return nil, err
	}
	token, err := s.authService.Login(ctx, in.GetEmail(), in.GetPassword())
	if err != nil {
		if errors.Is(err, customerr.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &pb.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Logout(ctx context.Context, in *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	return &pb.LogoutResponse{}, nil
}

func validateLogin(in *pb.LoginRequest) error {
	if in.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is empty")
	}

	if in.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is empty")
	}
	return nil
}

func validateRegister(in *pb.RegisterRequest) error {
	if in.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is empty")
	}

	if in.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is empty")
	}
	return nil
}
