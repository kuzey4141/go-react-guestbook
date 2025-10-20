package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// openDB (DEĞİŞİKLİK YOK)
func openDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// YENİ: CORS Middleware'i
// Bu fonksiyon, React (localhost:3000) sunucumuzdan gelen isteklere
// tarayıcının izin vermesini sağlayan başlıkları (headers) ekler.
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Sadece 'localhost:3000'den gelen isteklere izin ver
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		// İzin verilen metodlar
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		// İzin verilen başlıklar (özellikle POST'taki 'Content-Type' için)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Tarayıcılar, POST yapmadan önce bir 'OPTIONS' isteği (preflight) gönderir.
		// Bu isteğe 'Tamam, izin veriyorum' dememiz gerekir.
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Her şey tamamsa, asıl handler'a (mux) devam et
		next.ServeHTTP(w, r)
	})
}

func main() {
	// 1. Veritabanı Bağlantısı (DEĞİŞİKLİK YOK - Şifrenizin doğru olduğundan emin olun)
	connStr := "postgres://kuzey:sifreniz@localhost:5432/guestbook_db?sslmode=disable"
	db, err := openDB(connStr)
	if err != nil {
		log.Fatalf("Veritabanına bağlanılamadı: %v", err)
	}
	fmt.Println("Veritabanına başarıyla bağlanıldı.")

	// 2. Bağımlılıkları Hazırla
	// YENİ: 'templates' satırını sildik
	store := NewMessageStore(db)

	// 3. Application struct'ını oluştur
	// YENİ: 'templates' alanını sildik
	app := &Application{
		store: store,
	}

	// 4. Yönlendiriciyi (Router) ayarla
	mux := http.NewServeMux()
	// YENİ: Yönlendirme (route) artık '/' değil, '/api/messages'
	mux.HandleFunc("/api/messages", app.handleMessages)

	// 5. Web Sunucusunu yapılandır
	srv := &http.Server{
		Addr: ":8080",
		// YENİ: Mux'ı CORS middleware'i ile sarmaladık
		Handler: enableCORS(mux),
	}

	// 6. Sunucuyu Başlat
	fmt.Println("API Sunucusu http://localhost:8080 adresinde başlatılıyor...")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("Sunucu başlatılamadı:", err)
	}
}
