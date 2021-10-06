package httptrace

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dusted-go/diagnostic/log"
	"github.com/dusted-go/diagnostic/trace"
	"github.com/dusted-go/http/v2/middleware/chain"
	"github.com/dusted-go/http/v2/request"
)

// GetTraceFunc gets or generates trace IDs from an incoming HTTP request.
type GetTraceFunc func(r *http.Request, logger log.Event) (trace.ID, trace.SpanID)

// CreateLogEventFunc creates a new log event.
type CreateLogEventFunc func() log.Event

// Init is a middleware which initialised tracing information and a request scoped log event.
func Init(getTrace GetTraceFunc) func(CreateLogEventFunc) chain.Intermediate {
	return func(createLogEvent CreateLogEventFunc) chain.Intermediate {
		return func(next http.Handler) http.Handler {
			fn := func(w http.ResponseWriter, r *http.Request) {
				// Initialise a new log event for this request pipeline
				logger := createLogEvent().
					SetHTTPRequest(r).
					AddLabel("requestPath", r.URL.Path)

				// Get trace and span IDs for current request
				traceID, spanID := getTrace(r, logger)

				// Decorate log event with trace
				logger = logger.
					SetTraceID(traceID).
					SetSpanID(spanID)

				// Store trace and log event data in context
				ctx := log.Context(r.Context(), logger)
				ctx = trace.Context(ctx, traceID, spanID)

				// Log incoming request
				logger.SetData(struct {
					Headers map[string][]string
				}{
					Headers: r.Header,
				}).Fmt("%s %s %s", r.Proto, r.Method, request.FullURL(r))

				// Execute next middleware
				next.ServeHTTP(w, r.WithContext(ctx))
			}
			return http.HandlerFunc(fn)
		}
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
			return traceID, spanID, fmt.Errorf("failed to parse trace ID: %w", err)
		}

		return traceID, trace.DefaultGenerator.NewSpanID(), nil
	}

	// Trace ID and Span ID have been submitted
	if len(values) == 2 {
		traceID, err := trace.ParseID(values[0])
		if err != nil {
			return traceID, spanID, fmt.Errorf("failed to parse trace ID: %w", err)
		}

		// Remove the (optional) sampling parameter
		spanIDValue := strings.SplitN(values[1], ";", 2)[0]

		spanID, err = trace.ParseGoogleCloudSpanID(spanIDValue)
		if err != nil {
			return traceID, spanID, fmt.Errorf("failed to parse span ID: %w", err)
		}
		return traceID, spanID, nil
	}

	// Bad or no data
	return traceID, spanID, errors.New("invalid trace value in HTTP header")

}

// GoogleCloudTrace initialises tracing using the X-Cloud-Trace-Context HTTP header.
var GoogleCloudTrace = Init(
	func(r *http.Request, logger log.Event) (trace.ID, trace.SpanID) {
		traceHeader := r.Header.Get("X-Cloud-Trace-Context")
		if len(traceHeader) == 0 {
			return trace.DefaultGenerator.NewTraceIDs()
		}

		traceID, spanID, err := parseGoogleTraceContext(traceHeader)
		if err != nil {
			logger.Alert().SetError(err).Fmt("Invalid X-Cloud-Trace-Context header: %s", traceHeader)
			return trace.DefaultGenerator.NewTraceIDs()
		}

		return traceID, spanID
	})
