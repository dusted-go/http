package health

import (
	"fmt"
	"net/http"
)

func PingPong(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/ping" {
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprint(w, "pong")
				return
			}
			next.ServeHTTP(w, r)
		})
}
