package model

import "time"

type Message struct {
	ID        int64     `db:"id"`
	Username  string    `db:"username"`
	Text      string    `db:"text"`
	ChatID    int64     `db:"chat_id"`
	CreatedAt time.Time `db:"created_at"`
}
