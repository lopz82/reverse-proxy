package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type JSONError struct {
	Message string `json:"message"`
}

const ReverseProxyAddr = "0.0.0.0:9090"

var upServers string
var kind string

var serversPool []url.URL
var DefaultPath = "/:redirect"

func main() {

	flag.StringVar(&upServers, "upservers", "", "Upstream servers")
	flag.StringVar(&kind, "type", "", "Load balancing strategy")
	flag.Parse()

	r := gin.Default()

	for _, route := range config.Routes {

		servers, err := validateURLs(route.Servers)
		if err != nil {
			log.Panic(fmt.Sprintf("Error parsings urls: %s", err))
		}
		lb, err := createLoadBalancer(route.Strategy, servers, registeredFactories)
		if err != nil {
			log.Panic(fmt.Sprintf("Error creating a %s load balancer: %v", kind, err))
		}
		r.GET(DefaultPath, ReverseProxy(lb))
	}

	if err := r.Run(ReverseProxyAddr); err != nil {
		log.Printf("Error: %v", err)
	}
}

func ReverseProxy(lb LoadBalancer) gin.HandlerFunc {

	return func(c *gin.Context) {
		destination := lb.next()
		req := c.Request

		req.URL.Scheme = destination.Scheme
		req.URL.Host = destination.Host

		transport := http.DefaultTransport
		resp, err := transport.RoundTrip(req)
		if err != nil {
			msg := fmt.Sprintf("Server %s not available", req.URL.String())
			log.Println(msg)
			c.JSON(500, JSONError{
				Message: msg,
			})
			return
		}

		for k, gv := range resp.Header {
			for _, v := range gv {
				c.Header(k, v)
			}
		}
		defer resp.Body.Close()
		bufio.NewReader(resp.Body).WriteTo(c.Writer)
		return
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
