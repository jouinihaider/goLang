package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"estiam/dictionary"
	"estiam/route"
)

const dictionaryFile = "your_dictionary_file.gob"

func main() {
	d, err := dictionary.New(dictionaryFile)
	if err != nil {
		fmt.Println("Error creating dictionary:", err)
		return
	}

	// Set up Gorilla Mux router and register routes
	r := route.SetupRoutes(d)

	// Start the HTTP server concurrently
	go func() {
		port := 8080
		addr := fmt.Sprintf(":%d", port)
		fmt.Printf("Server listening on %s...\n", addr)
		http.ListenAndServe(addr, r)
	}()

	// Start the console-based interface
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Choose interaction method (console, http, exit):")
		method, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		method = strings.TrimSpace(method)

		switch method {
		case "console":
			consoleInterface(d, reader)
		case "http":
			// HTTP interface is already running in the goroutine
			fmt.Println("Use a tool like curl or a web browser to interact with the HTTP API.")
		case "exit":
			fmt.Println("Exiting program.")
			return
		default:
			fmt.Println("Invalid method. Please enter console, http, or exit.")
		}
	}
}

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
			actionAdd(d, reader)
		case "define":
			actionDefine(d, reader)
		case "remove":
			actionRemove(d, reader)
		case "list":
			actionList(d)
		case "back":
			return
		default:
			fmt.Println("Invalid action. Please enter add, define, remove, list, or back.")
		}
	}
}


func actionAdd(d *dictionary.Dictionary, reader *bufio.Reader) {
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

func actionDefine(d *dictionary.Dictionary, reader *bufio.Reader) {
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

func actionRemove(d *dictionary.Dictionary, reader *bufio.Reader) {
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

func actionList(d *dictionary.Dictionary) {
	words, entries := d.List()
	fmt.Println("Words in the dictionary:")
	for _, word := range words {
		fmt.Printf("%s: %s\n", word, entries[word].String())
	}
}
