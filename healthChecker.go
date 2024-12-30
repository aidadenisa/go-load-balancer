package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type HealthChecker struct {
	servers *[]Server
	healthCheckIntervalSeconds int
}

func NewHealthChecker(servers *[]Server, healthCheckIntervalSeconds int) HealthChecker{
	return HealthChecker{
		servers: servers,
		healthCheckIntervalSeconds: healthCheckIntervalSeconds,
	}
}

func (h *HealthChecker) checkServersHealth(wg *sync.WaitGroup, errChan chan<- error, ctx context.Context) {
	defer wg.Done()
	httpClient := http.Client{}
	
	for {
		select {
		case <-ctx.Done(): 
			fmt.Println("Stopping the health check")
			return 

		default: 
			for index, server := range *h.servers {

				request, err := http.NewRequest(http.MethodGet, server.URL + server.HealthCheckPath, nil)
				if err != nil {
					// server.Healthy = false
					errChan <- fmt.Errorf("error while creating the request for health check: %w", err)
					continue
				}

				response, err := httpClient.Do(request)
				if err != nil {
					(*h.servers)[index].Healthy = false
					continue
				}

				(*h.servers)[index].Healthy = response.StatusCode == http.StatusOK
			}
			time.Sleep(time.Duration(h.healthCheckIntervalSeconds) * time.Second)
		}
	}
}
