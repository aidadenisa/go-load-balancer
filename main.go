package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type RoundRobinBalancer struct {
	servers []string
	counter int
}

func NewRoundRobinBalancer() RoundRobinBalancer{
	return RoundRobinBalancer{
		servers: []string{"http://localhost:3200", "http://localhost:3201", "http://localhost:3202"},
		counter: -1,
	}
}

func (b *RoundRobinBalancer) nextServer() string {
	b.counter = b.counter + 1
	// next server is at the index of (the current counter) mod (the length of the servers list)
	return b.servers[b.counter % len(b.servers)]
}


func forwardReq(balancer RoundRobinBalancer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serverURL := balancer.nextServer()

		fmt.Printf("%s %s from %s to %s, user agent %s, accept %s\n", r.Method, r.URL.Path, r.RemoteAddr, serverURL, r.UserAgent(), r.Header.Get("Accept"))
		httpClient := http.Client{}

		// Create new req to hardcoded BE server
		request, err := http.NewRequest(r.Method, serverURL + r.URL.RequestURI(), r.Body)
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
		fmt.Printf("Resp %d, body: %s\n", resp.StatusCode, string(body))
		w.Write(body)
	}
}



func main() {
	balancer := NewRoundRobinBalancer()

	http.HandleFunc("/", forwardReq(balancer))

	log.Fatal(http.ListenAndServe(":3222", nil))
}