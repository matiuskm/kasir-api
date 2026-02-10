package middlewares

import (
	"log"
	"net/http"
	"time"
)

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf("[REQUEST] %s %s from %s", r.Method, r.RequestURI, r.RemoteAddr)
		
		next(w, r)

		duration := time.Since(start)
		log.Printf("[RESPONSE] %s %s completed in %v", r.Method, r.RequestURI, duration)
	}
}