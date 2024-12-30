package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)


func forwardReq(balancer RoundRobinBalancer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		server, err := balancer.NextServer()
		if err != nil {
			fmt.Println(err)
			http.Error(w, "error", 500)
			return 
		}

		httpClient := http.Client{}

		// Create new req to hardcoded BE server
		request, err := http.NewRequest(r.Method, server.URL + r.URL.RequestURI(), r.Body)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "error", 500)
			return 
		}

		// copy original headers
		for name, values := range r.Header {
			for _, value := range values {
				request.Header.Add(name, value)
			}
		}

		// make request
		resp, err := httpClient.Do(request)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "error", 500)
			return 
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "error", 500)
			return 
		}

		fmt.Printf("%s %s from %s to %s, user agent %s, accept %s\n", r.Method, r.URL.Path, r.RemoteAddr, server.URL, r.UserAgent(), r.Header.Get("Accept"))
		fmt.Printf("Resp %d, body: %s\n", resp.StatusCode, string(body))
		w.Write(body)
	}
}



func main() {
	healthCheckIntervalSeconds := flag.Int("healthCheckSec", 10, "health check interval in seconds")
	flag.Parse()

	servers := []Server{
		{ URL: "http://localhost:3200", Healthy: false, HealthCheckPath: "/health"}, 
		{ URL: "http://localhost:3201", Healthy: false, HealthCheckPath: "/health"}, 
		{ URL: "http://localhost:3202", Healthy: false, HealthCheckPath: "/health"}, 
	}

	errChan := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())
	

	// Health Check
	healthChecker := NewHealthChecker(&servers, *healthCheckIntervalSeconds)
	go healthChecker.checkServersHealth(errChan, ctx)

	// Round Robin Balancer
	balancer := NewRoundRobinBalancer(&servers)

	// Routes
	http.HandleFunc("/", forwardReq(balancer))

	// Serve the http server
	go func() {
		if err := http.ListenAndServe(":3222", nil) ; err != nil {
			errChan <- err
		}
	}()

	err := <-errChan
	cancel()
	log.Fatal(err)
}