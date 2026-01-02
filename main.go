package main

import (
	"go-postgres-stocks/router"
	"log"
	"log/slog"
	"net/http"
)

func main() {
	r := router.Router()

	slog.Info("starting the server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
