package request

import (
	"fmt"
	"net/http"
	"strings"
)

// IsHTTPS checks if a request object is HTTPS or not.
func IsHTTPS(r *http.Request) bool {
	return strings.HasPrefix(r.Proto, "HTTP") &&
		(r.TLS != nil || strings.ToLower(r.Header.Get("X-Forwarded-Proto")) == "https")
}

func fullURL(r *http.Request, forceScheme string) string {
	scheme := forceScheme
	if scheme == "" {
		if IsHTTPS(r) {
			scheme = "https"
		} else {
			scheme = "http"
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

// FullURL returns the request's full URL.
func FullURL(r *http.Request) string {
	return fullURL(r, "")
}

// HTTPSURL returns the request's full URL with https:// regardless of the original scheme.
func HTTPSURL(r *http.Request) string {
	return fullURL(r, "https")
}

// ShiftPath splits off the first component of path.
// Head will never contain a slash and tail will always be a rooted path without trailing slash.
func ShiftPath(path string) (head, tail string) {
	i := strings.Index(path[1:], "/") + 1
	if i <= 0 {
		return path[1:], "/"
	}
	return path[1:i], path[i:]
}
