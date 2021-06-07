package chain

import "net/http"

// IntermediateFunc is a middleware func which wraps another middleware func.
type IntermediateFunc func(http.HandlerFunc) http.HandlerFunc

// Intermediate is a middleware which wraps another middleware.
type Intermediate func(http.Handler) http.Handler

// Final is a middleware which wraps the final http handler func to be executed.
type Final func(http.HandlerFunc) http.Handler

// CreateFinal composes a bigger HTTP application from a series of middlewares.
// The middlewares are chained in ascending order: Create(A, B, C) will apply A(B(C)).
func CreateFinal(middlewares ...Intermediate) Final {
	return func(next http.HandlerFunc) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			final := http.Handler(next)
			for i := len(middlewares) - 1; i >= 0; i-- {
				m := middlewares[i]
				final = m(final)
			}
			final.ServeHTTP(w, r)
		})
	}
}
