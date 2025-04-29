package main

import (
	"Golang_balancer/api/handler"
	"Golang_balancer/db"
	"Golang_balancer/db/migrations"
	"Golang_balancer/internal/balancer"
	"Golang_balancer/internal/config"
	"Golang_balancer/internal/middleware"
	"Golang_balancer/internal/ratelimiter"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	path := "configs" + string(os.PathSeparator) + "config.json"
	pathBucket := "configs" + string(os.PathSeparator) + "bucketConfig.json"
	cfg, err := config.LoadConfig(path)
	if err != nil {
		log.Printf("error load config %s", err)
	}
	bucketCfg, err := config.LoadBucketConfig(pathBucket)
	if err != nil {
		log.Printf("error load bucket config %s", err)
	}

	bm := &ratelimiter.BucketManager{Config: bucketCfg}
	mux := http.NewServeMux()
	backendPool := &balancer.BackedPool{
		BackendsInfo: cfg.BackendsInfo,
	}
	h := handler.Handler{Pool: backendPool}
	mux.HandleFunc("/balancer", h.BalanceHandler)
	logMux := middleware.LogMiddleware(mux)
	bucketMux := bm.BucketMiddleware(logMux)

	server := &http.Server{
		Addr:         cfg.Port,
		Handler:      bucketMux,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}
	go backendPool.InitCheck()
	go backendPool.HealthCheck(5 * time.Second)
	log.Printf("Starting server at %s", server.Addr)

	dataBase, err := db.Connect()
	if err != nil {
		log.Fatalln("Failed connect to database")
	}
	migrations.Migrate(dataBase)
	err = server.ListenAndServe()
	if err != nil {
		log.Printf("error start server %s", err)
		return
	}

}
