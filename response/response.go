package response

import (
	"errors"
	"net/http"
	"syscall"

	"github.com/dusted-go/fault/fault"
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

	// If the error is a "broken pipe" then ignore it.
	// (this basically means the connection was aborted/closed by the peer)
	if errors.Is(err, syscall.EPIPE) {
		return nil
	}

	if err != nil {
		return fault.SystemWrapf(err, "response", "Plaintext",
			"Failed to write plaintext message to HTTP response body: %s.", text)
	}
	return nil
}
