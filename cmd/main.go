package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"applicationDesignTest/internal/api"
	"applicationDesignTest/internal/service/booking"
	"applicationDesignTest/internal/service/notify"
	"applicationDesignTest/internal/storage"
)

func main() {

	booker := booking.New(
		storage.New(),
		notify.New(),
	)

	cfg := LoadConfig()
	a := api.New(booker)

	go func() {
		err := a.Start(cfg.listen)
		if errors.Is(err, http.ErrServerClosed) {
			slog.Info("Server stopped")
		} else {
			slog.Error("Server failed: %s", err.Error())
			os.Exit(1)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
	slog.Info("got OS signal, stop application")

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	err := a.Stop(ctx)
	if err != nil {
		slog.Error("got error while stop: %s" + err.Error())
	}

}

type config struct {
	listen string
}

func LoadConfig() config {
	return config{
		listen: ":8080",
	}
}
