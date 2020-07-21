package main

import (
	"encoding/json"
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
	MESSAGE_COUNT_PREFIX = "m_c"
)

func queryMessagesWithUser(chanID, lastID int64) ([]Message, error) {
	return scanMessagesWithUser(db.Query("SELECT * FROM message INNER JOIN user ON user.id = message.user_id WHERE message.id > ? AND message.channel_id = ? ORDER BY message.id DESC LIMIT 100",
		lastID, chanID))
}

type MessageCount struct {
	Count     int64 `db:"count"`
	ChannelID int64 `db:"channel_id"`
}

func initMessageCount() error {
	var mc []MessageCount
	err := db.Select(&mc, "SELECT count(*) as count, channel_id FROM message GROUP BY channel_id")
	if err != nil {
		return err
	}
	for _, v := range mc {
		key := fmt.Sprintf("%s:%s", MESSAGE_COUNT_PREFIX, v.ChannelID)
		cache.Set(key, v.Count)
	}
	return nil
}

func IncrMessageCount(chID int64) error {
	key := fmt.Sprintf("%s:%s", MESSAGE_COUNT_PREFIX, chID)
	return cache.Increment(key, 1)
}

func FetchMessageCount(chID int64) (int64, error) {
	key := fmt.Sprintf("%s:%s", MESSAGE_COUNT_PREFIX, chID)
	data, err := cache.Get(key)
	if err != nil {
		return 0, err
	}
	var count int64
	err = json.Unmarshal(data, &count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
