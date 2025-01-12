package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func OpenTelemetryMiddleware(serviceName string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the tracer
			tracer := otel.Tracer(serviceName)

			// Start a new span for the incoming request
			ctx, span := tracer.Start(r.Context(), r.URL.Path)
			defer span.End()

			// Add attributes to the span
			span.SetAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.url", r.URL.String()),
				attribute.String("http.client_ip", r.RemoteAddr),
			)

			// Pass the context with the span to the next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
