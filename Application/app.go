package application

import {
	
	"fmt"
	"net/http"
	"github.com/go-chi/chi/v5"
	"context"
}

type App struct {
	router http.Handler
}



// build the constructor for the application

func New() *App {
	app := &App{
		router: loadRoutes()

		
	}
	return app
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:			":3000",
		Handler:	a.router,
	}
	err := server.ListenAndServer()
	if err != nil {
		return fmt.Errorf("Failed to start server: %w", err)
	}
	return nil
}



