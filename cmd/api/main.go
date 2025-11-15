package main

import (
	"fmt"
	"net/http"

	"github.com/thehappyidiot/save-my-music/internal/server"
)

func main() {
	smmServer := server.NewServer()

	fmt.Print("Whoop whoop (that's the sound of my server)")

	err := smmServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}
}
