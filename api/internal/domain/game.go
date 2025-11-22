package domain

type Game struct {
	ID        int64  `db:"id"`
	Name      string `db:"name"`
	OwnerID   int64  `db:"owner_id"`
	CreatedAt string `db:"created_at"`
	IsActive  bool   `db:"is_active"`
}

type Player struct {
	ID        int64  `db:"id"`
	GameID    int64  `db:"game_id"`
	UserID    int64  `db:"user_id"`
	JoinedAt  string `db:"joined_at"`
}
