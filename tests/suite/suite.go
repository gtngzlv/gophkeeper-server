package suite

import (
	"context"
	"net"
	"strconv"
	"testing"

	pb "github.com/gtngzlv/gophkeeper-protos/gen/go/gophkeeper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gtngzlv/gophkeeper-server/internal/config"
)

const grpcHost = "localhost"

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient pb.AuthClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper() // при фейле теста правильно формировался стек вызовов и эта функция не была указана как финальная
	t.Parallel()

	cfg := config.MustLoadByPath("../config/config.yaml")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)
	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	cc, err := grpc.DialContext(context.Background(),
		grpcAddress(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc connection failed: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: pb.NewAuthClient(cc),
	}
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))
}
