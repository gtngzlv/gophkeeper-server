package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

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

	go func() {
		application.GRPCSrv.MustRun(ctx)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGHUP)
	<-stop

	application.GRPCSrv.Stop(ctx)
	log.Info("Gracefully stopped")
}
