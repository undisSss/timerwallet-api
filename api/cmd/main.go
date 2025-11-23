package main

import (
	"log"
	"net/http"
	"os"

	"timerwallet-api/internal/db"
	"timerwallet-api/internal/delivery"

	"github.com/gorilla/mux"
)

func main() {
	connStr := os.Getenv("POSTGRES_CONN")
	db.InitDB(connStr)

	r := mux.NewRouter()

	gameHandler := &delivery.GameHandler{}

	// auth handlers moved to internal/delivery
	r.HandleFunc("/auth/telegram", delivery.TelegramAuthHandler).Methods("POST")
	r.HandleFunc("/auth/refresh", delivery.RefreshHandler).Methods("POST")
	r.HandleFunc("/auth/logout", delivery.LogoutHandler).Methods("POST")
	r.HandleFunc("/me", delivery.MeHandler).Methods("GET")

	r.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Data Exchange Success"))
	}).Methods("POST")

	r.HandleFunc("/games", gameHandler.CreateGame).Methods("POST")
	r.HandleFunc("/games", gameHandler.ListGames).Methods("GET")
	r.HandleFunc("/games/add-player", gameHandler.AddPlayer).Methods("POST")

	r.PathPrefix("/web/").Handler(http.StripPrefix("/web/", http.FileServer(http.Dir("./web"))))

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
