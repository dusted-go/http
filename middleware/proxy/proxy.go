package proxy

import (
	"net/http"
	"strings"

	"github.com/dusted-go/http/v3/middleware"
)

const https = "https"

// GetRealIP will take the first value which can be found in any of
// the given headers and then set the request's RemoteAddr with it.
func GetRealIP(headersToCheck ...string) middleware.Middleware {
	return middleware.Func(
		func(next http.Handler, w http.ResponseWriter, r *http.Request) {
			for _, headerName := range headersToCheck {
				val := r.Header.Get(headerName)
				if len(val) > 0 {
					r.RemoteAddr = val
					break
				}
			}
			next.ServeHTTP(w, r)
		},
	)
}

// ForwardedHeaders will parse the X-Forwarded-For and X-Forwarded-Proto
// headers and modify the request object accordingly.
//
// Set the proxyCount to the number of known proxies so that any values
// set by the origin caller get ignored.
func ForwardedHeaders(proxyCount int) middleware.Middleware {
	return middleware.Func(
		func(next http.Handler, w http.ResponseWriter, r *http.Request) {

			// Populate the URL object for later handlers
			r.URL.Scheme = "http"
			r.URL.Host = r.Host
			if r.TLS != nil && strings.HasPrefix(r.Proto, "HTTP") {
				r.URL.Scheme = https
			}

			// Skip if no proxies
			if proxyCount <= 0 {
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
			i := len(ips) - proxyCount
			if i < 0 {
				i = 0
			}
			realIP := strings.Trim(ips[i], " ")

			// Update request object
			r.RemoteAddr = realIP

			next.ServeHTTP(w, r)
		},
	)
}
