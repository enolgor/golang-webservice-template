package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/enolgor/golang-webservice-template/application"
	"github.com/enolgor/golang-webservice-template/config"
	"github.com/enolgor/golang-webservice-template/debug"
	"github.com/enolgor/golang-webservice-template/server"
	"github.com/enolgor/golang-webservice-template/shutdown"
)

var log *slog.Logger
var logLevel *slog.LevelVar

func init() {
	logLevel = new(slog.LevelVar)
	logLevel.Set(slog.LevelInfo)
	log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	config.Load(log, logLevel)
	if config.Mode == config.DEVELOPMENT && config.WaitForDebugger {
		debug.Wait(log)
	}
}

func main() {
	var app *application.App
	var mux *http.ServeMux
	var err error
	if app, err = application.New(log, logLevel); err != nil {
		log.Error("error starting application", "error", err.Error())
		os.Exit(-1)
	}
	if mux, err = server.NewMux(app); err != nil {
		log.Error("error starting application", "error", err.Error())
		os.Exit(-1)
	}
	srv := &http.Server{Addr: fmt.Sprintf(":%d", 8080), Handler: mux}
	go func() {
		log.Info("server started", "port", 8080)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Error("error starting server", "error", err.Error())
			os.Exit(-1)
		}
	}()
	shutdown := shutdown.NewShutdownGroup(log, 5*time.Second, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	shutdown.Register("server", func(ctx context.Context) error {
		return srv.Shutdown(ctx)
	})
	shutdown.Register("application", func(ctx context.Context) error {
		return app.Shutdown(ctx)
	})
	shutdown.WaitBlocking(context.Background())
}
