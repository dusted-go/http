package redirect

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dusted-go/http/middleware/chain"
	"github.com/dusted-go/http/request"
)

// ForceHTTPS is a middleware which redirects http:// requests to https://
func ForceHTTPS(host string, enable bool) chain.Intermediate {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if enable {
				redirect := r.Host == host && !request.IsHTTPS(r)
				if redirect {
					url := request.HTTPSURL(r)
					http.Redirect(w, r, url, 301)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

// TrailingSlash is a middleware which will redirect a matching request with
// a trailing slash in the path to the same endpoint without a trailing slash.
func TrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		redirect :=
			strings.HasPrefix(r.Proto, "HTTP") && // Must be HTTP request
				len(path) > 1 && // Skip if it is just the root (/) path
				path[len(path)-1] == '/' // Must have trailing slash

		if redirect {
			scheme := "http"
			if request.IsHTTPS(r) {
				scheme = "https"
			}
			url := fmt.Sprintf("%s://%s%s", scheme, r.Host, path[:len(path)-1])
			if r.URL.RawQuery != "" {
				url = url + "?" + r.URL.RawQuery
			}
			http.Redirect(w, r, url, 301)
			return
		}

		next.ServeHTTP(w, r)
	})
}
