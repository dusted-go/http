package response

import (
	"net/http"

	"github.com/dusted-go/utils/fault"
)

// ClearHeaders clears existing HTTP headers from the response.
func ClearHeaders(w http.ResponseWriter) {
	for k := range w.Header() {
		w.Header().Del(k)
	}
}

func Plaintext(
	statusCode int,
	text string,
	w http.ResponseWriter,
	r *http.Request,
) error {
	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", "text/plain")
	_, err := w.Write([]byte(text))
	if err != nil {
		return fault.SystemWrapf(err, "response", "Plaintext",
			"Failed to write plaintext message to HTTP response body: %s.", text)
	}
	return nil
}
