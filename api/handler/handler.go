package handler

import (
	"Golang_balancer/internal/balancer"
	"net/http"
	"net/http/httputil"
)

type Handler struct {
	pool *balancer.BackedPool
}

func (h *Handler) BalanceHandler(w http.ResponseWriter, r *http.Request) {
	target := h.pool.Next()
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ServeHTTP(w, r)
}
