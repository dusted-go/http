package recoverer

import (
	"errors"
	"net/http"
)

// RecoverFunc responds to a HTTP request which ended up panicking.
type RecoverFunc func(recovered any) http.HandlerFunc

// HandlePanics is a middleware which handles a panic and recovers gracefully by calling the RecovererFunc.
func HandlePanics(f RecoverFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				defer func() {
					if recovered := recover(); recovered != nil {
						if err, ok := recovered.(error); !ok || !errors.Is(err, http.ErrAbortHandler) {
							f(recovered)(w, r)
						}
					}
				}()
				next.ServeHTTP(w, r)
			},
		)
	}
}
