// dictionary.go

package dictionary

import (
	"errors"
	"sync"
)

type Entry struct {
	Definition string
}

func (e Entry) String() string {
	return e.Definition
}

type dictionaryOperation struct {
	word       string
	definition string
}

type Dictionary struct {
	entries   map[string]Entry
	operation chan dictionaryOperation
	mu        sync.RWMutex
}

func New() *Dictionary {
	d := &Dictionary{
		entries:   make(map[string]Entry),
		operation: make(chan dictionaryOperation),
	}
	go d.processOperations()
	return d
}

func (d *Dictionary) Add(word, definition string) {
	d.operation <- dictionaryOperation{word, definition}
}

func (d *Dictionary) Get(word string) (Entry, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	entry, found := d.entries[word]
	if !found {
		return Entry{}, errors.New("word not found in the dictionary")
	}
	return entry, nil
}

func (d *Dictionary) Remove(word string) {
	d.operation <- dictionaryOperation{word, ""} // Empty definition signifies removal
}

func (d *Dictionary) List() ([]string, map[string]Entry) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var words []string
	for word := range d.entries {
		words = append(words, word)
	}
	return words, d.entries
}

func (d *Dictionary) processOperations() {
	for {
		select {
		case op := <-d.operation:
			d.mu.Lock()
			if op.definition == "" {
				// If definition is empty, it's a removal operation
				delete(d.entries, op.word)
			} else {
				// Otherwise, it's an addition operation
				d.entries[op.word] = Entry{Definition: op.definition}
			}
			d.mu.Unlock()
		}
	}
}
