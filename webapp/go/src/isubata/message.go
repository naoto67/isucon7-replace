package main

import "time"

type Message struct {
	ID        int64     `db:"id"`
	ChannelID int64     `db:"channel_id"`
	UserID    int64     `db:"user_id"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`

	User *User
}

func queryMessagesWithUser(chanID, lastID int64) ([]Message, error) {
	return scanMessagesWithUser(db.Query("SELECT * FROM message INNER JOIN user ON user.id = message.user_id WHERE message.id > ? AND message.channel_id = ? ORDER BY message.id DESC LIMIT 100",
		lastID, chanID))
}
