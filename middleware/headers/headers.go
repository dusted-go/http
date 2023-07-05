package headers

import (
	"fmt"
	"net/http"

	"github.com/dusted-go/http/v4/middleware"
)

// Security will set HTTP security headers.
//
// It is the same as calling response.SecurityHeaders with the given hstsMaxAge, except that the middleware will apply those headers to all responses automatically.
//
// For more information check the response.SecurityHeaders function.
func Security(hstsMaxAge int) middleware.Middleware {
	return middleware.Func(
		func(next http.Handler, w http.ResponseWriter, r *http.Request) {
			w.Header().Set(
				"Strict-Transport-Security",
				fmt.Sprintf("max-age=%d; includeSubDomains", hstsMaxAge))
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "SAMEORIGIN")
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			next.ServeHTTP(w, r)
		},
	)
}
