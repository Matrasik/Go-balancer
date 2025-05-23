package handler

import (
	"Golang_balancer/internal/balancer"
	"log"
	"net/http"
	"net/http/httputil"
)

type Handler struct {
	Pool *balancer.BackedPool
}

// Основной хендлер, получает адресс бекенд-сервера и его id. Если не смог получить адрес - все сервера упали.
func (h *Handler) BalanceHandler(w http.ResponseWriter, r *http.Request) {
	target, n := h.Pool.Next()
	if target == nil {
		http.Error(w, "All servers are dead D:", http.StatusServiceUnavailable)
		log.Print("[Handler] All servers are dead!!!")
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		h.Pool.BackendsInfo[n].SetAlive(false) // Если нет коннекта - ставим серверу статус умер
		log.Printf("Error connect to backend server at Url: %s err: %v ", target.String(), err)
		http.Error(w, "Backend unreachable", http.StatusBadGateway)

	}
	proxy.ServeHTTP(w, r)
}
