package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
    "timerwallet-api/internal/db"
    "timerwallet-api/internal/delivery"
)

func main() {
	connStr := os.Getenv("POSTGRES_CONN")
	db.InitDB(connStr)

	r := mux.NewRouter()

	gameHandler := &delivery.GameHandler{} // TODO: внедрить usecase

	r.HandleFunc("/auth/telegram", TelegramAuthHandler).Methods("POST")
	r.HandleFunc("/data", DataExchangeHandler).Methods("POST")

	r.HandleFunc("/games", gameHandler.CreateGame).Methods("POST")
	r.HandleFunc("/games", gameHandler.ListGames).Methods("GET")
	r.HandleFunc("/games/add-player", gameHandler.AddPlayer).Methods("POST")

	r.PathPrefix("/web/").Handler(http.StripPrefix("/web/", http.FileServer(http.Dir("./web"))))

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func TelegramAuthHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement Telegram authentication logic
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Telegram Auth Success"))
}

func DataExchangeHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement data exchange logic
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data Exchange Success"))
}
