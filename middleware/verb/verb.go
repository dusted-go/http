package verb

import (
	"net/http"

	"github.com/dusted-go/http/middleware/chain"
)

func isAllowed(method string, methods []string) bool {
	for _, m := range methods {
		if m == method {
			return true
		}
	}
	return false
}

// Allow creates a middleware function which checks if an
// incoming HTTP request matches one of the allowed methods.
// If yes then it will invoke the next handler otherwise
// it will call the notAllowedFn function.
func Allow(methodNotAllowedFn http.HandlerFunc) func(...string) chain.IntermediateFunc {
	return func(methods ...string) chain.IntermediateFunc {
		return func(next http.HandlerFunc) http.HandlerFunc {
			fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if isAllowed(r.Method, methods) {
					next(w, r)
				} else {
					methodNotAllowedFn(w, r)
				}
			})
			return fn
		}
	}
}
