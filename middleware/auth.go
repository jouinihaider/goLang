// middleware/auth.go

package middleware

import (
    "net/http"
)

// AuthMiddleware is a middleware to authenticate requests based on a token.
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Get the token from the request header
        token := r.Header.Get("Authorization")

        // Validate the token (you should replace this logic with your actual authentication logic)
        if isValidToken(token) {
            // Token is valid, proceed to the next handler
            next.ServeHTTP(w, r)
        } else {
            // Token is invalid, respond with unauthorized status
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
        }
    })
}

// Replace this function with your actual token validation logic
func isValidToken(token string) bool {
    // Add your token validation logic here
    // For simplicity, this example considers any non-empty token as valid
    return token != ""
}
