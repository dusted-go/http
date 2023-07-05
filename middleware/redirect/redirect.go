package redirect

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dusted-go/http/v4/middleware"
)

const (
	httpsProto = "https"
	httpProto  = "http"
)

func isHTTPS(r *http.Request) bool {
	return strings.HasPrefix(r.Proto, "HTTP") &&
		(r.TLS != nil || strings.ToLower(r.Header.Get("X-Forwarded-Proto")) == httpsProto)
}

func fullURL(r *http.Request, desiredScheme string) string {
	scheme := desiredScheme
	if scheme == "" {
		if isHTTPS(r) {
			scheme = httpsProto
		} else {
			scheme = httpProto
		}
	}

	pathAndQuery := r.RequestURI
	if pathAndQuery == "" {
		pathAndQuery = r.URL.Path
		if r.URL.User.String() != "" {
			pathAndQuery = r.URL.User.String() + "@" + pathAndQuery
		}
		if r.URL.RawQuery != "" {
			pathAndQuery += "?" + r.URL.RawQuery
		}
		if r.URL.Fragment != "" {
			pathAndQuery += "#" + r.URL.Fragment
		}
	}

	return fmt.Sprintf("%s://%s%s", scheme, r.Host, pathAndQuery)
}

// ForceHTTPS is a middleware which redirects http:// requests to https://
func ForceHTTPS(enable bool, hosts ...string) middleware.Middleware {
	return middleware.Func(
		func(next http.Handler, w http.ResponseWriter, r *http.Request) {
			if enable {
				isMatch := false
				for _, h := range hosts {
					if h == r.Host {
						isMatch = true
						break
					}
				}
				redirect := isMatch && !isHTTPS(r)
				if redirect {
					url := fullURL(r, httpsProto)
					http.Redirect(w, r, url, 301)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
}

// TrailingSlash is a middleware which will redirect a matching request with
// a trailing slash in the path to the same endpoint without a trailing slash.
func TrailingSlash() middleware.Middleware {
	return middleware.Func(
		func(next http.Handler, w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			redirect :=
				strings.HasPrefix(r.Proto, "HTTP") && // Must be HTTP request
					len(path) > 1 && // Skip if it is just the root (/) path
					path[len(path)-1] == '/' // Must have trailing slash

			if redirect {
				scheme := httpProto
				if isHTTPS(r) {
					scheme = httpsProto
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
func Hosts(hosts map[string]string, enable bool) middleware.Middleware {
	return middleware.Func(
		func(next http.Handler, w http.ResponseWriter, r *http.Request) {
			if enable {
				if dest, ok := hosts[r.Host]; ok {
					r.Host = dest
					url := fullURL(r, "")
					http.Redirect(w, r, url, 301)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
}
