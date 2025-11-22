package delivery

import (
	"encoding/json"
	"net/http"
	"strconv"

	"timerwallet-api/internal/usecase"
)

type GameHandler struct {
	Usecase usecase.GameUsecase
}

func (h *GameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OwnerID int64  `json:"owner_id"`
		Name    string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request"))
		return
	}
	game, err := h.Usecase.CreateGame(req.OwnerID, req.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error creating game"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(game)
}

func (h *GameHandler) ListGames(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	archivedStr := r.URL.Query().Get("archived")
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)
	archived := archivedStr == "true"
	games, err := h.Usecase.ListGames(userID, archived)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error listing games"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(games)
}

func (h *GameHandler) AddPlayer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		GameID int64 `json:"game_id"`
		UserID int64 `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request"))
		return
	}
	if err := h.Usecase.AddPlayer(req.GameID, req.UserID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error adding player"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Player added"))
}
