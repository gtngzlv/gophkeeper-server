package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gtngzlv/gophkeeper-server/internal/proto/pb"

	"github.com/gtngzlv/gophkeeper-server/internal/app"
	"github.com/gtngzlv/gophkeeper-server/internal/config"
	"github.com/gtngzlv/gophkeeper-server/internal/logger"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background())
	defer cancel()
	cfg := config.MustLoad()
	log := logger.MustSetup(cfg.Env)

	application, err := app.NewApp(ctx, log, cfg)
	if err != nil {
		panic("failed to init application" + err.Error())
	}

	go application.GRPCSrv.MustRun()
	go runRest(cfg)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGHUP,
		syscall.SIGKILL,
		syscall.SIGSEGV)
	<-stop

	application.GRPCSrv.Stop(ctx)
	log.Info("Gracefully stopped")
}

func runRest(cfg *config.Config) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterGophkeeperHandlerFromEndpoint(ctx, mux, "localhost:"+string(cfg.GRPC.Port), opts)
	if err != nil {
		panic(err)
	}
	mux.HandlePath("GET", "/gophkeeper.swagger.json", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		http.FileServer(http.Dir("internal/proto")).ServeHTTP(w, r)
	})
	log.Printf("rest listening on port %v", cfg.REST.Port)
	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", cfg.REST.Port), mux); err != nil {
		panic(err)
	}
}
