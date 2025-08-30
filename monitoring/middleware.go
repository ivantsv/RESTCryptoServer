package monitoring

import (
	"net/http"
	"time"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"strconv"
)

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		ActiveConnections.Inc()
		defer ActiveConnections.Dec()
		
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		
		next.ServeHTTP(ww, r)
		
		duration := time.Since(start)
		status := strconv.Itoa(ww.Status())
		
		rctx := chi.RouteContext(r.Context())
		routePattern := rctx.RoutePattern()
		if routePattern == "" {
			routePattern = "not_found"
		}
		
		RecordHTTPRequest(r.Method, routePattern, status, duration)
		
		userID := r.Context().Value("user_id")
		if userID == nil {
			userID = "anonymous"
		}
		LogRequest(r.Method, r.URL.Path, ww.Status(), duration, userID.(string))
	})
}

func TracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := r.Header.Get("X-Trace-ID")
		if traceID == "" {
			traceID = generateTraceID()
		}
		
		ctx := Logger.With().Str("trace_id", traceID).Logger().WithContext(r.Context())
		
		w.Header().Set("X-Trace-ID", traceID)
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func generateTraceID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}