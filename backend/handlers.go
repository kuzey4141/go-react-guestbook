package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Application struct {
	store *MessageStore
}

func (app *Application) handleMessages(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		app.handleGetMessages(w, r)
	case "POST":
		app.handlePostMessage(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (app *Application) handleGetMessages(w http.ResponseWriter, r *http.Request) {
	messages, err := app.store.GetAll()
	if err != nil {
		log.Printf("database error (GetAll): %v", err)
		http.Error(w, "Could not load messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func (app *Application) handlePostMessage(w http.ResponseWriter, r *http.Request) {
	var newMessage Message

	err := json.NewDecoder(r.Body).Decode(&newMessage)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if newMessage.Name == "" || newMessage.Message == "" {
		http.Error(w, "Name and message are required", http.StatusBadRequest)
		return
	}

	if err := app.store.Add(newMessage); err != nil {
		log.Printf("database error (Add): %v", err)
		http.Error(w, "Could not save message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newMessage)
}
