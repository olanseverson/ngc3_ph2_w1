package config

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func SetupServer() (*httprouter.Router, *http.Server) {
	router := httprouter.New()

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}
	fmt.Println("Server running on port 8080")
	return router, &server
}
