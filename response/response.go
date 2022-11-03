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

// SecureHeaders sets the following HTTP headers for better security:
//
// - Strict-Transport-Security
//
// - X-Content-Type-Options
//
// - X-Frame-Options
//
// - Referrer-Policy
//
// For the Strict-Transport-Security you must set the maxAge parameter in seconds.
//
// See more: https://developer.mozilla.org/en-US/docs/Glossary/HSTS
func SecurityHeaders(
	w http.ResponseWriter,
	hstsMaxAge int,
) {
	w.Header().Set(
		"Strict-Transport-Security",
		fmt.Sprintf("max-age=%d; includeSubDomains", hstsMaxAge))
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
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
