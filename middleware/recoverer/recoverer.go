package recoverer

import (
	"errors"
	"net/http"

	"github.com/dusted-go/fault/stack"
	"github.com/dusted-go/http/v3/middleware"
)

// RecoverFunc responds to a HTTP request which ended up panicking.
type RecoverFunc func(recovered interface{}, stack stack.Trace) http.HandlerFunc

// HandlePanics is a middleware which handles a panic and recovers gracefully by calling the RecovererFunc.
func HandlePanics(f RecoverFunc) middleware.Middleware {
	return middleware.Func(
		func(next http.Handler, w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recovered := recover(); recovered != nil {
					if err, ok := recovered.(error); !ok || !errors.Is(err, http.ErrAbortHandler) {
						stackTrace := stack.Capture()
						f(recovered, *stackTrace)(w, r)
					}
				}
			}()
			next.ServeHTTP(w, r)
		},
	)
}
