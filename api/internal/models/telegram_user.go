package models

type TelegramUser struct {
	ID         int64  `db:"id"`
	TelegramID int64  `db:"telegram_id"`
	Username   string `db:"username"`
	FirstName  string `db:"first_name"`
	LastName   string `db:"last_name"`
	AuthDate   string `db:"auth_date"`
	CreatedAt  string `db:"created_at"`
}
