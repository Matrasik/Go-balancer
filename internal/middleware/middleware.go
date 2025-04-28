package middleware

import (
	"log"
	"net/http"
)

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			log.Printf("Request from addr: %s", r.RemoteAddr)
		}()
		next.ServeHTTP(w, r)
	})
}
