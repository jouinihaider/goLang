// middleware/logging.go

package middleware

import (
    "log"
    "net/http"
    "time"
)

// LoggingMiddleware logs information about incoming requests.
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Log the request information
        log.Printf("[%s] %s %s %s", time.Now().Format("2006-01-02 15:04:05"), r.Method, r.RequestURI, r.RemoteAddr)

        // Call the next handler in the chain
        next.ServeHTTP(w, r)
    })
}
