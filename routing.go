package main

import (
	"fmt"
	"log"
	"net/http"
)

func registerRoutes(mux *http.ServeMux) {
	for routeName, route := range config.Routes {

		servers, err := validateURLs(route.Servers)
		if err != nil {
			log.Panic(fmt.Sprintf("Error parsings urls: %s", err))
		}
		lb, err := createLoadBalancer(route.Strategy, servers, registeredFactories)
		if err != nil {
			log.Panic(fmt.Sprintf("Error creating a %s load balancer: %v", route.Strategy, err))
		}
		mux.HandleFunc(routeName, ReverseProxy(lb, route.Root))
		log.Printf("Route %s registered successfully", routeName)
	}
}
