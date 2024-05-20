package main

import (
	"fmt"
	"net/http"
)

func main() {

	server := &http.Server{
		Addr:    ":3000",
		Handler: http.HandlerFunc(basicHandler),
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
