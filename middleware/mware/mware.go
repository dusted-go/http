package mware

import "net/http"

func Bind(middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			middleware := middlewares[i]
			next = middleware(next)
		}
		return next
	}
}
