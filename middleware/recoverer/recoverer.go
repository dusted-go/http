package recoverer

import (
	"net/http"
	"runtime/debug"

	"github.com/dusted-go/http/middleware/chain"
)

// RecoverFunc responds to a HTTP request which ended up panicking.
type RecoverFunc func(rcv interface{}, stack []byte) http.HandlerFunc

// HandlePanics is a middleware which handles a panic and recovers gracefully by calling the RecovererFunc.
func HandlePanics(f RecoverFunc) chain.Intermediate {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rcv := recover(); rcv != nil && rcv != http.ErrAbortHandler {
					f(rcv, debug.Stack())(w, r)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
