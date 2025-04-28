package main

import (
	"Golang_balancer/api/handler"
	"Golang_balancer/internal/balancer"
	"Golang_balancer/internal/config"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	path := "configs" + string(os.PathSeparator) + "config.json"
	cfg, err := config.LoadConfig(path)
	if err != nil {
		log.Printf("error load config %s", err)
	}
	mux := http.NewServeMux()
	backendPool := &balancer.BackedPool{
		BackendsInfo: cfg.BackendsInfo,
	}
	h := handler.Handler{Pool: backendPool}
	mux.HandleFunc("/balancer", h.BalanceHandler)
	server := &http.Server{
		Addr:         cfg.Port,
		Handler:      mux,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}
	go backendPool.HealthCheck(5 * time.Second)
	err = server.ListenAndServe()
	if err != nil {
		log.Printf("error start server %s", err)
		return
	}

}
