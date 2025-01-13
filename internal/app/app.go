// Package app provides functionality to initialize, start, and stop the gophKeeper application.
// It sets up the GRPC server, repository, and database connections based on the configuration.
package app

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"mb-feedback/internal/api"
	"mb-feedback/internal/conf"
	"os"
	"os/signal"
)

// App represent the application state containing configuration, GRPC server, database connection, and repository.
type App struct {
	pgpool *pgxpool.Pool

	api *api.Rest

	exitCode int
}

// Init initializes the application with the provided configuration
func (a *App) Init() {
	var err error

	// pgpool
	{
		a.pgpool, err = pgxpool.New(context.Background(), conf.Conf.PgDsn)
		errCheck(err, "pgxpool.New")
	}

}

// Start starts the application, initializing and running GRPC server.
func (a *App) Start() {
	slog.Info("Starting")

}

// Listen listens for signals to stop the application
func (a *App) Listen() {
	signalCtx, signalCtxCancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer signalCtxCancel()

	// wait signal
	<-signalCtx.Done()
}

// Stop stops the application, shutting down the GRPC server.
func (a *App) Stop() {
	slog.Info("Shutting down...")

}

// Exit gracefully shuts down the application by logging the exit action
// and then terminating the program with the specified exit code.
func (a *App) Exit() {
	slog.Info("Exit")

	os.Exit(a.exitCode)
}

// errCheck checks if an error occurred and logs it with the specified message.
// If an error is found, the function logs the error and terminates the program.
// If a message is provided, it is included in the logged output.
func errCheck(err error, msg string) {
	if err != nil {
		if msg != "" {
			err = fmt.Errorf("%s: %w", msg, err)
		}
		slog.Error(err.Error())
		os.Exit(1)
	}
}
