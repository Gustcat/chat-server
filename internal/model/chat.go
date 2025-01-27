package model

import "time"

type Message struct {
	ID        int64     `json:"id"`
	Username  string    `json:"from"`
	Text      string    `json:"text"`
	ChatID    int64     `json:"chat_id"`
	CreatedAt time.Time `json:"timestamp"`
}
