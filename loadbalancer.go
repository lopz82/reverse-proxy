package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/url"
)

func init() {
	register("roundrobin", NewRoundRobinLoadBalancer, registeredFactories)
	register("random", NewRandomLoadBalancer, registeredFactories)
}

type LoadBalancerFactory func(pool []url.URL) LoadBalancer

type LoadBalancer interface {
	next() url.URL
}

type RoundRobinLoadBalancer struct {
	Idx  int
	Pool []url.URL
}

func (lb *RoundRobinLoadBalancer) next() url.URL {
	var res url.URL
	if lb.Idx == len(lb.Pool)-1 {
		res = lb.Pool[lb.Idx]
		lb.Idx = 0
		return res
	}
	res = lb.Pool[lb.Idx]
	lb.Idx += 1
	return res
}

func NewRoundRobinLoadBalancer(pool []url.URL) LoadBalancer {
	return &RoundRobinLoadBalancer{Pool: pool}
}

type RandomLoadBalancer struct {
	Pool []url.URL
}

func (lb *RandomLoadBalancer) next() url.URL {
	i := rand.Intn(len(lb.Pool))
	return lb.Pool[i]
}

func NewRandomLoadBalancer(pool []url.URL) LoadBalancer {
	return &RandomLoadBalancer{Pool: pool}
}

var registeredFactories = make(map[string]LoadBalancerFactory)

func register(name string, factory LoadBalancerFactory, mapping map[string]LoadBalancerFactory) {
	if name == "" {
		log.Panic("Cannot register factory with an empty name")
	}
	if factory == nil {
		log.Panic("Cannot register a empty factory")
	}
	_, registered := mapping[name]
	if registered {
		log.Printf("%s Factory already registered", name)
		return
	}
	mapping[name] = factory
}

func createLoadBalancer(kind string, pool []url.URL, mapping map[string]LoadBalancerFactory) (LoadBalancer, error) {
	factory, ok := mapping[kind]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Loadbalancer type %s not registered", kind))
	}
	return factory(pool), nil
}
