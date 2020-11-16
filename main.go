package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type JSONError struct {
	Message string `json:"message"`
}

const ReverseProxyAddr = "0.0.0.0:9090"

var upServers string
var kind string

var serversPool []url.URL
var DefaultPath = "/redirect"

func main() {

	flag.StringVar(&upServers, "upservers", "", "Upstream servers")
	flag.StringVar(&kind, "type", "", "Load balancing strategy")
	flag.Parse()

	mux := http.NewServeMux()

	for _, route := range config.Routes {

		servers, err := validateURLs(route.Servers)
		if err != nil {
			log.Panic(fmt.Sprintf("Error parsings urls: %s", err))
		}
		lb, err := createLoadBalancer(route.Strategy, servers, registeredFactories)
		if err != nil {
			log.Panic(fmt.Sprintf("Error creating a %s load balancer: %v", kind, err))
		}
		mux.HandleFunc(DefaultPath, ReverseProxy(lb))
	}

	s := &http.Server{
		Addr:           ReverseProxyAddr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Printf("Error: %v", err)
	}
}

func parseServersFlag(f string) []url.URL {
	servers := strings.Split(f, ",")
	var result []url.URL
	for _, s := range servers {

		addr, err := url.Parse(s)
		if err != nil {
			log.Printf("Error parsing %s: %v", s, err)
			continue
		}
		result = append(result, *addr)
	}
	return result
}
