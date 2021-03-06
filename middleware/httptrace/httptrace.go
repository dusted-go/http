package httptrace

import (
	"net/http"
	"strings"

	"github.com/dusted-go/diagnostic/v2/log"
	"github.com/dusted-go/diagnostic/v2/trace"
	"github.com/dusted-go/http/v3/request"
	"github.com/dusted-go/http/v3/server"
	"github.com/dusted-go/utils/fault"
)

// GetTraceFunc gets or generates trace IDs from an incoming HTTP request.
type GetTraceFunc func(r *http.Request) (trace.ID, trace.SpanID)

// CreateLogProviderFunc creates a new default log provider.
type CreateLogProviderFunc func() *log.Provider

// Init is a middleware which initialised tracing information and a request scoped log event.
func Init(getTrace GetTraceFunc) func(CreateLogProviderFunc) server.Middleware {
	return func(createLogProvider CreateLogProviderFunc) server.Middleware {
		return server.MiddlewareFunc(
			func(next http.Handler, w http.ResponseWriter, r *http.Request) {
				// Get a new log provider for this request and set and additional
				// label for the request path and the request itself:
				provider := createLogProvider()
				provider.AddLabel("requestPath", r.URL.Path)
				provider.SetHTTPRequest(r)

				// Update the request's context with the new log provider:
				r = r.WithContext(
					log.Context(r.Context(),
						provider))

				// Get trace and span IDs for current request and
				// store them in the request's context:
				traceID, spanID := getTrace(r)
				r = r.WithContext(
					trace.Context(r.Context(), traceID, spanID))

				// Log incoming request
				log.New(r.Context()).
					Data("requestHeaders", r.Header).
					Fmt("%s %s %s", r.Proto, r.Method, request.FullURL(r))

				// Execute next middleware
				next.ServeHTTP(w, r)
			},
		)
	}
}

func parseGoogleTraceContext(headerValue string) (trace.ID, trace.SpanID, error) {
	traceID := trace.ID{}
	spanID := trace.SpanID{}
	values := strings.SplitN(headerValue, "/", 2)

	// Only trace ID has been submitted
	if len(values) == 1 {
		// Remove the (optional) sampling parameter
		traceIDValue := strings.SplitN(values[0], ";", 2)[0]
		traceID, err := trace.ParseID(traceIDValue)
		if err != nil {
			return traceID, spanID,
				fault.SystemWrap(
					err, "httptrace", "parseGoogleTraceContext", "failed to parse traceID")
		}

		return traceID, trace.DefaultGenerator.NewSpanID(), nil
	}

	// Trace ID and Span ID have been submitted
	if len(values) == 2 {
		traceID, err := trace.ParseID(values[0])
		if err != nil {
			return traceID, spanID,
				fault.SystemWrap(
					err, "httptrace", "parseGoogleTraceContext", "failed to parse traceID")
		}

		// Remove the (optional) sampling parameter
		spanIDValue := strings.SplitN(values[1], ";", 2)[0]

		spanID, err = trace.ParseGoogleCloudSpanID(spanIDValue)
		if err != nil {
			return traceID, spanID,
				fault.SystemWrap(
					err, "httptrace", "parseGoogleTraceContext", "failed to parse spanID")
		}
		return traceID, spanID, nil
	}

	// Bad or no data
	return traceID, spanID,
		fault.System("httptrace", "parseGoogleTraceContext", "invalid trace value in HTTP header")
}

// GoogleCloudTrace initialises tracing using the X-Cloud-Trace-Context HTTP header.
var GoogleCloudTrace = Init(
	func(r *http.Request) (trace.ID, trace.SpanID) {
		traceHeader := r.Header.Get("X-Cloud-Trace-Context")
		if len(traceHeader) == 0 {
			return trace.DefaultGenerator.NewTraceIDs()
		}

		traceID, spanID, err := parseGoogleTraceContext(traceHeader)
		if err != nil {
			log.New(r.Context()).
				Alert().
				Err(err).
				Fmt("Invalid X-Cloud-Trace-Context header: %s", traceHeader)
			return trace.DefaultGenerator.NewTraceIDs()
		}

		return traceID, spanID
	})
