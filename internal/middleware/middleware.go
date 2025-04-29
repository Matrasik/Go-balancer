package middleware

import (
	"log"
	"net/http"
)

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			log.Println("Request from addr: ", r.RemoteAddr)
		}()
		next.ServeHTTP(w, r)
	})
}
