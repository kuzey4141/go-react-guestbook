package main

import (
	"encoding/json" // YENİ: JSON ile çalışmak için bu paketi ekledik
	"log"
	"net/http"
	// "html/template" paketini sildik
)

// YENİ: Application struct'ından 'templates' alanını kaldırdık
type Application struct {
	store *MessageStore
}

// YENİ: Ana handler'ımızın adını 'handleMessages' olarak değiştirdik
func (app *Application) handleMessages(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		app.handleGetMessages(w, r)
	case "POST":
		app.handlePostMessage(w, r)
	default:
		// Sadece GET ve POST'a izin veriyoruz
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// YENİ: handleGetMessages - Tüm mesajları JSON olarak döndürür
func (app *Application) handleGetMessages(w http.ResponseWriter, r *http.Request) {
	mesajlar, err := app.store.GetAll()
	if err != nil {
		log.Printf("Veritabanı hatası (GetAll): %v", err)
		http.Error(w, "Mesajlar yüklenemedi", http.StatusInternalServerError)
		return
	}

	// YENİ: JSON yanıtı gönderme
	// 1. Başlık (Header) bilgisini "bu bir JSON'dur" olarak ayarla
	w.Header().Set("Content-Type", "application/json")
	// 2. 'mesajlar' dizisini JSON formatına çevir ve yanıt olarak gönder
	json.NewEncoder(w).Encode(mesajlar)
}

// YENİ: handlePostMessage - Yeni bir mesajı JSON olarak alır ve kaydeder
func (app *Application) handlePostMessage(w http.ResponseWriter, r *http.Request) {
	var yeniMesaj Mesaj

	// YENİ: Gelen isteğin gövdesindeki (body) JSON'u oku
	// ve 'yeniMesaj' struct'ına ata
	err := json.NewDecoder(r.Body).Decode(&yeniMesaj)
	if err != nil {
		http.Error(w, "Geçersiz JSON formatı", http.StatusBadRequest)
		return
	}

	// Form doğrulaması (basitçe)
	if yeniMesaj.Isim == "" || yeniMesaj.Mesaj == "" {
		http.Error(w, "İsim ve Mesaj alanları boş olamaz", http.StatusBadRequest)
		return
	}

	// Veritabanına ekle
	if err := app.store.Add(yeniMesaj); err != nil {
		log.Printf("Veritabanı hatası (Add): %v", err)
		http.Error(w, "Mesaj kaydedilemedi", http.StatusInternalServerError)
		return
	}

	// YENİ: Başarılı yanıtı gönder
	// 201 Created (Oluşturuldu) durum koduyla
	// ve kaydedilen mesajı JSON olarak geri döndür
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(yeniMesaj)
}
