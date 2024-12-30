package main

import "fmt"

type Server struct {
	URL string
	Healthy bool
	HealthCheckPath string
}

type RoundRobinBalancer struct {
	servers *[]Server
	counter int
}

func NewRoundRobinBalancer(servers *[]Server) RoundRobinBalancer{
	balancer := RoundRobinBalancer{
		servers: servers,
		counter: -1,
	}

	return balancer
}

func (b *RoundRobinBalancer) NextServer() (Server, error) {
	countUnhealthy := 0
	
	var next Server
	// If all the nodes are unhealthy, stop
	for countUnhealthy < len(*b.servers){
		// Check the next node
		b.counter = b.counter + 1

		// next server is at the index of (the current counter) mod (the length of the servers list)
		next = (*b.servers)[b.counter % len(*b.servers)]

		if !next.Healthy {
			countUnhealthy++
		} else {
			break
		}
	}

	if countUnhealthy == len(*b.servers) {
		return Server{},fmt.Errorf("all nodes are unhealthy")
	}

	return next, nil
}