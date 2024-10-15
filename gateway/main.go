package main

import (
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	"github.com/karokojnr/bookmesh-gateway/handler"
	shared "github.com/karokojnr/bookmesh-shared"
)

var (
	httpAddr = shared.EnvString("HTTP_ADDR", ":8080")
)

func main() {
	mux := http.NewServeMux()
	h := handler.NewHandler()
	h.RegisterRoutes(mux)

	log.Println("Starting http server on", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("Failed to start http server", err)
	}
}
