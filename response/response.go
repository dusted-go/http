package response

import "net/http"

// ClearHeaders clears existing HTTP headers from the response.
func ClearHeaders(w http.ResponseWriter) {
	for k := range w.Header() {
		w.Header().Del(k)
	}
}
