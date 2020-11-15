package main

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestRegister(t *testing.T) {
	testCasesRegister := []struct {
		Name    string
		Factory LoadBalancerFactory
		Success bool
		Panics  bool
	}{
		{"", NewRoundRobinLoadBalancer, false, true},
		{"FakeName", nil, false, true},
		{"NewLoadBalancer", NewRandomLoadBalancer, true, false},
	}
	for _, test := range testCasesRegister {
		registeredFactories := make(map[string]LoadBalancerFactory)
		if test.Panics {
			assert.Panics(t, func() {
				register(test.Name, test.Factory, registeredFactories)
			})
			assert.True(t, len(registeredFactories) == 0)
			continue
		}
		if test.Success {
			// Registering should be idempotent
			register(test.Name, test.Factory, registeredFactories)
			register(test.Name, test.Factory, registeredFactories)
			assert.True(t, len(registeredFactories) == 1)
		} else {
			register(test.Name, test.Factory, registeredFactories)
			assert.True(t, len(registeredFactories) == 0)
		}
	}
}

func TestCreateLoadBalancer(t *testing.T) {
	testCasesCreateLoadBalancer := []struct {
		Kind     string
		Expected interface{}
		Error    bool
	}{
		{"roundrobin", &RoundRobinLoadBalancer{}, false},
		{"random", &RandomLoadBalancer{}, false},
		{"nonexistent", nil, true},
	}

	registry := map[string]LoadBalancerFactory{
		"roundrobin": NewRoundRobinLoadBalancer,
		"random":     NewRandomLoadBalancer,
	}

	for _, test := range testCasesCreateLoadBalancer {
		u, _ := url.Parse("http://testing.com")
		pool := []url.URL{*u}
		lb, err := createLoadBalancer(test.Kind, pool, registry)
		if test.Error {
			assert.Error(t, err)
			continue
		}
		assert.IsType(t, test.Expected, lb)
		assert.True(t, lb.next() == *u)
	}
}

var u1, _ = url.Parse("http://testing.com")
var u2, _ = url.Parse("http://testing2.com")
var pool = []url.URL{*u1, *u2}

func TestRoundRobinLoadBalancer(t *testing.T) {
	lb := NewRoundRobinLoadBalancer(pool)
	assert.True(t, lb.next() == *u1)
	assert.True(t, lb.next() == *u2)
	assert.True(t, lb.next() == *u1)
}

func TestRandomLoadBalancer(t *testing.T) {
	lb := NewRandomLoadBalancer(pool)
	assert.Contains(t, pool, lb.next())
	assert.Contains(t, pool, lb.next())
	assert.Contains(t, pool, lb.next())
}
