package main

import (
	"fmt"
	"time"
)

type Message struct {
	ID        int64     `db:"id"`
	ChannelID int64     `db:"channel_id"`
	UserID    int64     `db:"user_id"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`

	User *User
}

const (
	MESSAGE_SET = "m_set"
)

func queryMessagesWithUser(chanID, lastID int64) ([]Message, error) {
	return scanMessagesWithUser(db.Query("SELECT * FROM message INNER JOIN user ON user.id = message.user_id WHERE message.id > ? AND message.channel_id = ? ORDER BY message.id DESC LIMIT 100",
		lastID, chanID))
}

func addMessage(channelID, userID int64, content string) (int64, error) {
	now := time.Now()
	res, err := db.Exec(
		"INSERT INTO message (channel_id, user_id, content, created_at) VALUES (?, ?, ?, ?)",
		channelID, userID, content, now)
	if err != nil {
		return 0, err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	err = AddMessageCache(Message{
		ID:        lastID,
		ChannelID: channelID,
		UserID:    userID,
		Content:   content,
		CreatedAt: now,
	})
	return lastID, err
}

func initMessagesCache() {
	var msgs []Message
	db.Select(&msgs, "SELECT * FROM message")
	for _, v := range msgs {
		AddMessageCache(v)
	}
	return
}

func AddMessageCache(m Message) error {
	key := fmt.Sprintf("%s:%s", MESSAGE_SET, m.ChannelID)
	return cache.ZAdd(key, m.ID, m)
}

func fetchMessageCount(chID, lastID int64) (int64, error) {
	key := fmt.Sprintf("%s:%s", MESSAGE_SET)
	return cache.ZCount(key, lastID+1, "+inf")
}
