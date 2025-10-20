package main

import (
	"database/sql"
	"fmt"
)

// YENİ: Struct'ımıza JSON etiketleri (tags) ekledik.
// `json:"isim"` etiketi, Go'daki 'Isim' alanının
// JSON'a 'isim' olarak (küçük harfle) çevrilmesini sağlar.
type Mesaj struct {
	Isim  string `json:"isim"`
	Mesaj string `json:"mesaj"`
}

// MessageStore (DEĞİŞİKLİK YOK)
type MessageStore struct {
	db *sql.DB
}

// NewMessageStore (DEĞİŞİKLİK YOK)
func NewMessageStore(db *sql.DB) *MessageStore {
	return &MessageStore{
		db: db,
	}
}

// Add (DEĞİŞİKLİK YOK)
func (s *MessageStore) Add(msg Mesaj) error {
	sqlStatement := `INSERT INTO messages (name, message) VALUES ($1, $2)`
	_, err := s.db.Exec(sqlStatement, msg.Isim, msg.Mesaj)
	if err != nil {
		return fmt.Errorf("mesaj eklenemedi: %w", err)
	}
	return nil
}

// GetAll (DEĞİŞİKLİK YOK)
func (s *MessageStore) GetAll() ([]Mesaj, error) {
	sqlStatement := `SELECT name, message FROM messages ORDER BY created_at DESC`

	rows, err := s.db.Query(sqlStatement)
	if err != nil {
		return nil, fmt.Errorf("mesajlar alınamadı: %w", err)
	}
	defer rows.Close()

	var mesajlar []Mesaj
	for rows.Next() {
		var msg Mesaj
		// .Scan'in sırası SELECT ile aynı olmalı (name, message -> &msg.Isim, &msg.Mesaj)
		if err := rows.Scan(&msg.Isim, &msg.Mesaj); err != nil {
			return nil, fmt.Errorf("satır taranamadı: %w", err)
		}
		mesajlar = append(mesajlar, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("satır okuma hatası: %w", err)
	}

	return mesajlar, nil
}
