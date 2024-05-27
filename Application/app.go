package application

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

type App struct {
	router http.Handler
	rdb    *redis.Client
	config Config
}

// build the constructor for the application

func New(config Config) *App {
	app := &App{
		rdb: redis.NewClient(&redis.Options{
			Addr: config.RedisAddress,
		}),
		config: config,
	}
	app.loadRoutes()
	return app
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.config.ServerPort),
		Handler: a.router,
	}
	// see if redis is there
	err := a.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("Failed to connect to redis: %w", err)
	}

	defer func() {
		if err := a.rdb.Close(); err != nil {
			fmt.Println("Failed to close redis", err)
		}
	}()

	fmt.Println("Starting server")

	ch := make(chan error, 1)

	go func() {

		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("Failed to start server: %w", err)
		}

		close(ch)
	}()

	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		return server.Shutdown(timeout)
	}
	return nil

}
