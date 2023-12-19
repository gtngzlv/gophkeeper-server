package keeper

import (
	"context"

	pb "github.com/gtngzlv/gophkeeper-protos/gen/go/gophkeeper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gtngzlv/gophkeeper-server/internal/domain/models"
)

type IKeeperService interface {
	SaveData(ctx context.Context, data []models.KeeperData) error
}

type serverAPI struct {
	pb.UnimplementedKeeperServer

	keeper IKeeperService
}

func Register(grpcServer *grpc.Server, k IKeeperService) {
	pb.RegisterKeeperServer(grpcServer, &serverAPI{keeper: k})
}

func (s *serverAPI) SaveData(ctx context.Context, in *pb.SaveDataRequest) (*pb.SaveDataResponse, error) {
	data, err := pbDataToDomain(in)
	if err != nil {
		return nil, err
	}

	if err = s.keeper.SaveData(ctx, data); err != nil {
		return nil, status.Error(codes.Internal, "failed to save data")
	}
	return &pb.SaveDataResponse{}, nil
}

func pbDataToDomain(in *pb.SaveDataRequest) ([]models.KeeperData, error) {
	var data []models.KeeperData
	for _, v := range in.GetData() {
		if v == "" {
			return nil, status.Error(codes.InvalidArgument, "data is empty")
		}
		md := models.KeeperData{Value: v}
		data = append(data, md)
	}
	return data, nil
}
