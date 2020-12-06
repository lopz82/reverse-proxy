package main

import (
	"log"
	"net/http"
	"time"
)

func main() {

	mux := http.NewServeMux()
	registerRoutes(mux)

	s := &http.Server{
		Addr:           config.ProxyConfig.ProxyAddress,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Reverse proxy listening on %s", config.ProxyConfig.ProxyAddress)
	if err := s.ListenAndServe(); err != nil {
		log.Printf("Error: %v", err)
	}
}

