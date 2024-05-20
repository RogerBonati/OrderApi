package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Get("/hello", basicHandler)

	server := &http.Server{
		Addr:    ":3000",
		Handler: router,
	}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Failed to listen to server: ", err)
	}

}

func basicHandler(w http.ResponseWriter, r *http.Request) {
	// cast inputstring to bytearray
	_, err := w.Write([]byte("Servus Sepp!"))
	if err != nil {
		fmt.Println("Server could not write message: ", err)
	}
}
