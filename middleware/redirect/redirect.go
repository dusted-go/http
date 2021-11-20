package redirect

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dusted-go/http/v2/request"
	"github.com/dusted-go/http/v2/server"
)

// ForceHTTPS is a middleware which redirects http:// requests to https://
func ForceHTTPS(enable bool, hosts ...string) server.Middleware {
	return server.MiddlewareFunc(
		func(next http.Handler, w http.ResponseWriter, r *http.Request) {
			if enable {
				isMatch := false
				for _, h := range hosts {
					if h == r.Host {
						isMatch = true
						break
					}
				}
				redirect := isMatch && !request.IsHTTPS(r)
				if redirect {
					url := request.HTTPSURL(r)
					http.Redirect(w, r, url, 301)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
}

// TrailingSlash is a middleware which will redirect a matching request with
// a trailing slash in the path to the same endpoint without a trailing slash.
func TrailingSlash() server.Middleware {
	return server.MiddlewareFunc(
		func(next http.Handler, w http.ResponseWriter, r *http.Request) {
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

// Hosts is a middleware which redirects from one host to another.
// (e.g. https://www.foo.bar -> https://foo.bar)
func Hosts(hosts map[string]string, enable bool) server.Middleware {
	return server.MiddlewareFunc(
		func(next http.Handler, w http.ResponseWriter, r *http.Request) {
			if enable {
				if dest, ok := hosts[r.Host]; ok {
					r.Host = dest
					url := request.FullURL(r)
					http.Redirect(w, r, url, 301)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
}
