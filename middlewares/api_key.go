package middlewares

import "net/http"

func APIKey(validApiKey string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// before function running
			apiKey := r.Header.Get("X-API-Key")
			if apiKey == "" {
				http.Error(w, "API key is missing", http.StatusUnauthorized)
				return
			}

			if apiKey != validApiKey {
				http.Error(w, "Invalid API key", http.StatusForbidden)
				return
			}

			next(w, r)

			// after function running
		}
	}
}