package main

import (
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/appengine/memcache"
)

type HaveRead struct {
	UserID    int64     `db:"user_id"`
	ChannelID int64     `db:"channel_id"`
	MessageID int64     `db:"message_id"`
	UpdatedAt time.Time `db:"updated_at"`
	CreatedAt time.Time `db:"created_at"`
}

const (
	HAVEREAD_PREFIX = "hr"
)

func AddHaveReadCache(haveread HaveRead) error {
	key := fmt.Sprintf("%s:%s:%s", HAVEREAD_PREFIX, haveread.UserID, haveread.ChannelID)
	return cache.Set(key, haveread)
}

func FetchHaveReadCache(userID, chID int64) (*HaveRead, error) {
	key := fmt.Sprintf("%s:%s:%s", HAVEREAD_PREFIX, userID, chID)
	data, err := cache.Get(key)
	if err != nil {
		return nil, err
	}
	var haveread HaveRead
	err = json.Unmarshal(data, &haveread)
	return &haveread, err
}

func queryHaveRead(userID, chID int64) (int64, error) {
	haveread, err := FetchHaveReadCache(userID, chID)

	if err != memcache.ErrCacheMiss {
		return 0, nil
	}
	return haveread.MessageID, nil
}
