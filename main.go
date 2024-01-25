// main.go

package main

import (
    "bufio"
    "fmt"
    "net/http"
    "strings"
    "encoding/json"  // Add this line for the json package

    "github.com/gorilla/mux"

    "estiam/dictionary"
    "estiam/route"
    "estiam/middleware"
)

const dictionaryFile = "your_dictionary_file.gob"

func main() {
    d, err := dictionary.New(dictionaryFile)
    if err != nil {
        fmt.Println("Error creating dictionary:", err)
        return
    }

    r := route.SetupRoutes(d)

    // Start the HTTP server with the router
    http.Handle("/", r)

    // Use the logging middleware globally for all routes
    http.ListenAndServe(":8080", middleware.LoggingMiddleware(http.DefaultServeMux))
}

// Existing functions for console-based interactions

func consoleInterface(d *dictionary.Dictionary, reader *bufio.Reader) {
    for {
        fmt.Println("Enter action (add, define, remove, list, back):")
        action, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("Error reading input:", err)
            continue
        }

        action = strings.TrimSpace(action)

        switch action {
        case "add":
            consoleActionAdd(d, reader)
        case "define":
            consoleActionDefine(d, reader)
        case "remove":
            consoleActionRemove(d, reader)
        case "list":
            consoleActionList(d)
        case "back":
            return
        default:
            fmt.Println("Invalid action. Please enter add, define, remove, list, or back.")
        }
    }
}

func consoleActionAdd(d *dictionary.Dictionary, reader *bufio.Reader) {
    fmt.Println("Enter word:")
    word, err := reader.ReadString('\n')
    if err != nil {
        fmt.Println("Error reading input:", err)
        return
    }
    word = strings.TrimSpace(word)

    fmt.Println("Enter definition:")
    definition, err := reader.ReadString('\n')
    if err != nil {
        fmt.Println("Error reading input:", err)
        return
    }
    definition = strings.TrimSpace(definition)

    d.Add(word, definition)
    fmt.Printf("Word '%s' added to the dictionary.\n", word)
}

func consoleActionDefine(d *dictionary.Dictionary, reader *bufio.Reader) {
    fmt.Println("Enter word:")
    word, err := reader.ReadString('\n')
    if err != nil {
        fmt.Println("Error reading input:", err)
        return
    }
    word = strings.TrimSpace(word)

    entry, err := d.Get(word)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Printf("Definition of '%s': %s\n", word, entry.String())
}

func consoleActionRemove(d *dictionary.Dictionary, reader *bufio.Reader) {
    fmt.Println("Enter word:")
    word, err := reader.ReadString('\n')
    if err != nil {
        fmt.Println("Error reading input:", err)
        return
    }
    word = strings.TrimSpace(word)

    d.Remove(word)
    fmt.Printf("Word '%s' removed from the dictionary.\n", word)
}

func consoleActionList(d *dictionary.Dictionary) {
    words, entries := d.List()
    fmt.Println("Words in the dictionary:")
    for _, word := range words {
        fmt.Printf("%s: %s\n", word, entries[word].String())
    }
}

// Functions for HTTP-based interactions


func httpActionAdd(d *dictionary.Dictionary, w http.ResponseWriter, r *http.Request) {
	var entry dictionary.Entry

	// Decode the request body into the Entry struct
	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Add the entry to the dictionary
	d.Add(entry.Word, entry.Definition)

	// Respond with success message
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Entry added to the dictionary"))
}

func httpActionDefine(d *dictionary.Dictionary, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	word := vars["word"]

	// Get the entry from the dictionary
	entry, err := d.Get(word)
	if err != nil {
		http.Error(w, "Word not found", http.StatusNotFound)
		return
	}

	// Respond with the definition
	response := map[string]string{"word": word, "definition": entry.String()}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func httpActionRemove(d *dictionary.Dictionary, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	word := vars["word"]

	// Remove the entry from the dictionary
	d.Remove(word)

	// Respond with success message
	w.Write([]byte("Entry removed from the dictionary"))
}

func httpActionList(d *dictionary.Dictionary, w http.ResponseWriter, r *http.Request) {
	// Get the list of words and entries from the dictionary
	words, entries := d.List()

	// Prepare the response as a JSON object
	response := make(map[string]map[string]string)
	for _, word := range words {
		response[word] = map[string]string{"definition": entries[word].String()}
	}

	// Convert the response to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
