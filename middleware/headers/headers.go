package headers

import (
	"net/http"

	"github.com/dusted-go/http/v3/middleware"
	"github.com/dusted-go/http/v3/response"
)

// Security will set HTTP security headers.
//
// It is the same as calling response.SecurityHeaders with the given hstsMaxAge, except that the middleware will apply those headers to all responses automatically.
//
// For more information check the response.SecurityHeaders function.
func Security(hstsMaxAge int) middleware.Middleware {
	return middleware.Func(
		func(next http.Handler, w http.ResponseWriter, r *http.Request) {
			response.SecurityHeaders(w, hstsMaxAge)
			next.ServeHTTP(w, r)
		},
	)
}
