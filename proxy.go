package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func ReverseProxy(lb LoadBalancer, root bool) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		prepareRedirection(lb.next(), r)

		transport := http.DefaultTransport

		if root == true {
			r.URL.Path = "/"
		}

		resp, err := transport.RoundTrip(r)
		if err != nil {
			msg := fmt.Sprintf("Server %s not available", r.URL.String())
			log.Println(msg)
			jsonResp, err := json.Marshal(JSONError{Message: msg})
			if err != nil {
				log.Fatal("Error formating 500 JSONError Response")
			}
			w.Header().Add("content-type", "application/json")
			w.WriteHeader(http.StatusBadGateway)
			_, err = w.Write(jsonResp)
			if err != nil {
				log.Fatal("Error writing 500 JSONError Response")
			}
			return
		}

		for k, gv := range resp.Header {
			for _, v := range gv {
				r.Header.Add(k, v)
			}
		}
		defer resp.Body.Close()
		bufio.NewReader(resp.Body).WriteTo(w)
		log.Printf("%s \"%s %s %s\" %d", r.RemoteAddr, r.Method, r.RequestURI, r.Proto, resp.StatusCode)
		return
	}
}

func prepareRedirection(dest url.URL, r *http.Request) {
	r.URL.Scheme = dest.Scheme
	r.URL.Host = dest.Host
}
