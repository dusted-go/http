package firewall

import (
	"net"
	"net/http"
	"strings"
)

func LimitRequestSize(
	maxSize int64,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				r.Body = http.MaxBytesReader(w, r.Body, maxSize)
				next.ServeHTTP(w, r)
			})
	}
}

func RestrictByIP(
	whitelisted []net.IP,
	unauthorized http.HandlerFunc,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if len(whitelisted) == 0 {
					next.ServeHTTP(w, r)
					return
				}
				requestIP := net.ParseIP(strings.Split(r.RemoteAddr, ":")[0])
				if requestIP == nil {
					unauthorized(w, r)
					return
				}
				for _, ip := range whitelisted {
					if ip.Equal(requestIP) {
						next.ServeHTTP(w, r)
						return
					}
				}
				unauthorized(w, r)
			})
	}
}
