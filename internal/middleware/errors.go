package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/murkh/gig-mobile-backend/internal/platform/logger"
	"go.uber.org/zap"
)

func Recoverer(l *logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					logEntry := middleware.GetLogEntry(r)
					if logEntry != nil {
						logEntry.Panic(rvr, debug.Stack())
					} else {
						l.Error("Panic recovered", zap.Any("panic", rvr), zap.ByteString("stack", debug.Stack()))
					}

					w.WriteHeader(http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// ErrorHandler is a middleware that handles errors returned by handlers
// Note: In Chi, handlers usually write to the response writer directly.
// This middleware is more useful if we had a custom handler type that returns errors.
// For standard http.Handler, we rely on the handler to log errors or use a wrapper.
// Here we will provide a wrapper function for handlers.
