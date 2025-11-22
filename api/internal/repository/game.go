package repository

import "timerwallet-api/internal/domain"

type GameRepository interface {
	Create(game *domain.Game) error
	FindByID(id int64) (*domain.Game, error)
	ListByUser(userID int64, archived bool) ([]domain.Game, error)
}

type PlayerRepository interface {
	Add(player *domain.Player) error
	ListByGame(gameID int64) ([]domain.Player, error)
}
