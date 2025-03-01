package social

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/usmartpro/social-network/internal/config"
	"github.com/usmartpro/social-network/internal/logger"
)

func main() {
	configuration, err := config.LoadConfiguration()
	if err != nil {
		log.Fatalf("Error read configuration: %s", err)
	}
	logg, err := logger.New(configuration.Logger)
	if err != nil {
		log.Println("error create logger: " + err.Error())
		os.Exit(1)
	}

	storageConf := storage.New(ctx, configuration.Storage.Dsn).Connect(ctx)
	socialNetwork := app.New(logg, storageConf)

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
}
