package gophkeeper

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	customerr "github.com/gtngzlv/gophkeeper-server/internal/domain/errors"
	"github.com/gtngzlv/gophkeeper-server/internal/domain/models"
	"github.com/gtngzlv/gophkeeper-server/internal/proto/pb"
)

type IGophkeeperService interface {
	Register(ctx context.Context, email string, password string) (userID int64, err error)
	Login(ctx context.Context, email string, password string) (token string, err error)
	SaveData(ctx context.Context, data models.PersonalData) error
}

type serverAPI struct {
	pb.UnimplementedGophkeeperServer

	service IGophkeeperService
}

func Register(grpcServer *grpc.Server, srv IGophkeeperService) {
	pb.RegisterGophkeeperServer(grpcServer, &serverAPI{
		service: srv,
	})
}

func (s *serverAPI) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if err := validateRegister(in); err != nil {
		return nil, err
	}

	userID, err := s.service.Register(ctx, in.GetEmail(), in.GetPassword())
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
	token, err := s.service.Login(ctx, in.GetEmail(), in.GetPassword())
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

func (s *serverAPI) SaveData(ctx context.Context, in *pb.SaveDataRequest) (*pb.SaveDataResponse, error) {
	if len(in.Data) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	data, err := pbDataToDomain(in)
	if err != nil {
		return nil, err
	}

	if err = s.service.SaveData(ctx, data); err != nil {
		if errors.Is(err, customerr.ErrFailedGetUserID) {
			return nil, status.Error(codes.Unauthenticated, "not logged in")
		}
		return nil, status.Error(codes.Internal, "failed to save data")
	}
	return &pb.SaveDataResponse{}, nil
}

func pbDataToDomain(in *pb.SaveDataRequest) (models.PersonalData, error) {
	var data []models.Data
	for _, v := range in.GetData() {
		if v == "" {
			return models.PersonalData{}, status.Error(codes.InvalidArgument, "data is empty")
		}
		md := models.Data{Value: v}
		data = append(data, md)
	}
	return models.PersonalData{PData: data}, nil
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
