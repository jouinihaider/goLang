// route/router.go

package route

import (
    "estiam/dictionary"
    "net/http"

    "github.com/gorilla/mux"
    "estiam/middleware"
)

// SetupRoutes initializes and returns the Gorilla Mux router with routes.
func SetupRoutes(d *dictionary.Dictionary) *mux.Router {
    r := mux.NewRouter()

    // Add the authentication middleware to secure your routes
    r.Use(middleware.AuthMiddleware)

    // Add your existing routes
    r.HandleFunc("/add", AddEntryHandler(d)).Methods("POST")
    r.HandleFunc("/get/{word}", GetDefinitionHandler(d)).Methods("GET")
    r.HandleFunc("/remove/{word}", RemoveEntryHandler(d)).Methods("DELETE")

    return r
}

func AddEntryHandler(d *dictionary.Dictionary) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Your logic to handle adding an entry using d.Add(...)
        // Example: Read data from request body, parse JSON, and call d.Add(...)
    }
}

func GetDefinitionHandler(d *dictionary.Dictionary) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Your logic to handle getting a definition using d.Get(...)
        // Example: Extract word from URL path, call d.Get(...) and return response
    }
}

func RemoveEntryHandler(d *dictionary.Dictionary) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Your logic to handle removing an entry using d.RemoveByWord(...)
        // Example: Extract word from URL path and call d.RemoveByWord(...)
    }
}
