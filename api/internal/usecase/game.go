package usecase

import "timerwallet-api/internal/domain"

type GameUsecase interface {
	CreateGame(ownerID int64, name string) (*domain.Game, error)
	ListGames(userID int64, archived bool) ([]domain.Game, error)
	AddPlayer(gameID, userID int64) error
}
