package proxy

import (
	"net/http"
	"strings"

	"github.com/dusted-go/http/v2/middleware/chain"
)

const https = "https"

// GetRealIP will parse the X-Forwarded-For header for the real remote address.
func GetRealIP(count int) chain.Intermediate {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Populate the URL object for later handlers
			r.URL.Scheme = "http"
			r.URL.Host = r.Host
			if r.TLS != nil && strings.HasPrefix(r.Proto, "HTTP") {
				r.URL.Scheme = https
			}

			// Skip if no proxies
			if count <= 0 {
				next.ServeHTTP(w, r)
				return
			}

			proto := r.Header.Get("X-Forwarded-Proto")
			if strings.ToLower(proto) == https {
				r.URL.Scheme = https
			}

			// Skip if no IP addresses set
			h := r.Header.Get("X-Forwarded-For")
			if len(h) == 0 {
				next.ServeHTTP(w, r)
				return
			}

			// Find the correct IP address of the client
			ips := strings.Split(h, ",")
			i := len(ips) - count
			if i < 0 {
				i = 0
			}
			realIP := strings.Trim(ips[i], " ")

			// Update request object
			r.RemoteAddr = realIP

			next.ServeHTTP(w, r)
		})
	}
}
