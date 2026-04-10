package main

import (
	"database/sql"
	"fmt"
)

// Message represents a guestbook entry and defines its JSON shape.
type Message struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

type MessageStore struct {
	db *sql.DB
}

func NewMessageStore(db *sql.DB) *MessageStore {
	return &MessageStore{
		db: db,
	}
}

func (s *MessageStore) Add(msg Message) error {
	sqlStatement := `INSERT INTO messages (name, message) VALUES ($1, $2)`
	_, err := s.db.Exec(sqlStatement, msg.Name, msg.Message)
	if err != nil {
		return fmt.Errorf("could not add message: %w", err)
	}
	return nil
}

func (s *MessageStore) GetAll() ([]Message, error) {
	sqlStatement := `SELECT name, message FROM messages ORDER BY created_at DESC`

	rows, err := s.db.Query(sqlStatement)
	if err != nil {
		return nil, fmt.Errorf("could not fetch messages: %w", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		// The Scan order must match the SELECT order.
		if err := rows.Scan(&msg.Name, &msg.Message); err != nil {
			return nil, fmt.Errorf("could not scan row: %w", err)
		}
		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration failed: %w", err)
	}

	return messages, nil
}
