package application

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type App struct {
	router http.Handler
}

// build the constructor for the application

func New() *App {
	app := &App{
		router: loadRoutes(),
	}
	return app
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":3000",
		Handler: a.router,
	}
	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("Failed to start server: %w", err)
	}
	return nil
}
