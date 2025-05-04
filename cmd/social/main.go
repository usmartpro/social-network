package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"social-network/internal/app"
	"social-network/internal/cache"
	"social-network/internal/config"
	"social-network/internal/logger"
	"social-network/internal/storage"
	"syscall"
	"time"

	internalhttp "social-network/internal/server/http"
)

func main() {
	var logg *logger.Logger
	configuration, err := config.LoadConfiguration()
	if err != nil {
		log.Fatalf("Error read configuration: %s", err)
	}
	logg, err = logger.New(configuration.Logger)
	if err != nil {
		log.Println("error create logger: " + err.Error())
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	storageConf := storage.New(ctx, configuration.Storage.Dsn).Connect(ctx)

	cacheConf := cache.New(ctx, configuration.Cache.Dsn).Connect(ctx)

	socialNetwork := app.New(logg, storageConf, cacheConf)

	// HTTP
	server := internalhttp.NewServer(logg, socialNetwork, configuration.HTTP.Host, configuration.HTTP.Port)

	go func() {
		if err := server.Start(ctx); err != nil {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
		}
	}()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("social network is running...")

	<-ctx.Done()
}
