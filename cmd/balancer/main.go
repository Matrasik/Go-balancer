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
	// Загрузка конфигов
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
	// Подключение к бд и миграции
	dataBase, err := db.Connect()
	if err != nil {
		log.Fatalln("Failed connect to database")
	}
	migrations.Migrate(dataBase)

	// создание роутов и сервера
	bm := ratelimiter.NewBucketManager(*bucketCfg, dataBase)
	mux := http.NewServeMux()
	backendPool := &balancer.BackedPool{
		BackendsInfo: cfg.BackendsInfo,
	}
	h := handler.Handler{Pool: backendPool}
	mux.HandleFunc("/", h.BalanceHandler)
	bucketMux := bm.BucketMiddleware(mux)
	logMux := middleware.LogMiddleware(bucketMux)
	server := &http.Server{
		Addr:         cfg.Port,
		Handler:      logMux,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	// Запуск сервера и HealtCheck'ов бекендов. Один сразу же, другой каждые пять секунд.
	go backendPool.InitCheck()
	go backendPool.HealthCheck(5 * time.Second)
	log.Printf("Starting server at %s", server.Addr)

	err = server.ListenAndServe()
	if err != nil {
		log.Printf("error start server %s", err)
		return
	}

}
