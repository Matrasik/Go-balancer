package main

import (
	"fmt"
	"net/http"
)

func runServer(addr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			fmt.Fprintln(w, "Addr:", addr, "URL:", r.URL.String())
		})
	mux.HandleFunc("/health",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
		})
	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	fmt.Println("starting server at", addr)
	server.ListenAndServe()
}

func main() {
	go runServer(":8081")
	go runServer(":8082")
	go runServer(":8083")
	go runServer(":8084")
	go runServer(":8085")
	go runServer(":8086")
	go runServer(":8087")
	go runServer(":8088")
	go runServer(":8089")
	runServer(":8090")
}
