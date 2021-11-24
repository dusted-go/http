package recoverer

import (
	"errors"
	"net/http"
	"runtime/debug"

	"github.com/dusted-go/http/v3/server"
)

// RecoverFunc responds to a HTTP request which ended up panicking.
type RecoverFunc func(recovered interface{}, stack []byte) http.HandlerFunc

// HandlePanics is a middleware which handles a panic and recovers gracefully by calling the RecovererFunc.
func HandlePanics(f RecoverFunc) server.Middleware {
	return server.MiddlewareFunc(
		func(next http.Handler, w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recovered := recover(); recovered != nil {
					if err, ok := recovered.(error); !ok || !errors.Is(err, http.ErrAbortHandler) {
						f(recovered, debug.Stack())(w, r)
					}
				}
			}()
			next.ServeHTTP(w, r)
		},
	)
}
