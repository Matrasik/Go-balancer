package main

import (
	"log"
	"net"
	"net/http"
	"time"
)

// Да, основной код читает конфиг (этим кодом как раз ставил испытания), но проверить как-либо лучше сложновато на локалхосте 😥

func main() {
	localIP := "localhost"
	localPort := 55554

	localAddr := &net.TCPAddr{
		IP:   net.ParseIP(localIP),
		Port: localPort,
	}

	dialer := &net.Dialer{
		LocalAddr: localAddr,
		Timeout:   10 * time.Second,
	}

	transport := &http.Transport{
		DialContext: dialer.DialContext,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	resp, err := client.Get("http://localhost:8080/")
	if err != nil {
		log.Println("Request error: :", err)
		return
	}
	defer resp.Body.Close()

}
