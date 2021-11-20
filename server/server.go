package server

import "net/http"

// Middleware allows to chain multiple HTTP handlers together.
type Middleware interface {
	Next(http.Handler) http.HandlerFunc
}

// MiddlewareFunc implements the Middleware interface on a function
// of type func(http.Handler, http.ResponseWriter, *http.Request).
type MiddlewareFunc func(http.Handler, http.ResponseWriter, *http.Request)

func (f MiddlewareFunc) Next(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f(next, w, r)
	}
}

// CombineMiddlewares chains one or many middlewares into a single middleware.
func CombineMiddlewares(middlewares ...Middleware) Middleware {
	return MiddlewareFunc(
		func(next http.Handler, w http.ResponseWriter, r *http.Request) {
			for i := len(middlewares) - 1; i >= 0; i-- {
				m := middlewares[i]
				next = m.Next(next)
			}
			next.ServeHTTP(w, r)
		},
	)
}
