package keeper

import (
	"context"
	"errors"

	pb "github.com/gtngzlv/gophkeeper-protos/gen/go/gophkeeper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	customerr "github.com/gtngzlv/gophkeeper-server/internal/domain/errors"
	"github.com/gtngzlv/gophkeeper-server/internal/domain/models"
)

type IKeeperService interface {
	SaveData(ctx context.Context, data models.PersonalData) error
}

type serverAPI struct {
	pb.UnimplementedKeeperServer

	keeper IKeeperService
}

func Register(grpcServer *grpc.Server, k IKeeperService) {
	pb.RegisterKeeperServer(grpcServer, &serverAPI{keeper: k})
}

func (s *serverAPI) SaveData(ctx context.Context, in *pb.SaveDataRequest) (*pb.SaveDataResponse, error) {
	if len(in.Data) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	data, err := pbDataToDomain(in)
	if err != nil {
		return nil, err
	}

	if err = s.keeper.SaveData(ctx, data); err != nil {
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
