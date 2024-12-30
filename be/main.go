package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func handleReq(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[BE] %s %s from %s, user agent %s, accept %s\n", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), r.Header.Get("Accept"))
	w.Write([]byte("Hello from the BE server!\n"))
}

func main() {
	portFlag := flag.Int("port", 3200, "port to run the server on")
	flag.Parse()

	http.HandleFunc("/", handleReq)

	fmt.Printf("Server is running on port %d\n", *portFlag)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *portFlag), nil))
}