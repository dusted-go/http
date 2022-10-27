package response

import (
	"fmt"
	"net/http"
)

// ClearHeaders clears existing HTTP headers from the response.
func ClearHeaders(w http.ResponseWriter) {
	for k := range w.Header() {
		w.Header().Del(k)
	}
}

// WritePlaintext writes a text/plaintext message to the response stream.
func WritePlaintext(
	w http.ResponseWriter,
	statusCode int,
	text string,
) error {
	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", "text/plain")
	_, err := w.Write([]byte(text))

	if err != nil {
		return fmt.Errorf("failed to respond with plaintext message '%s': %w", text, err)
	}
	return nil
}
