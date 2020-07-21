package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
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

func AddHaveRead(haveread HaveRead) error {
	_, err := db.Exec("INSERT INTO haveread (user_id, channel_id, message_id, updated_at, created_at)"+
		" VALUES (?, ?, ?, NOW(), NOW())"+
		" ON DUPLICATE KEY UPDATE message_id = ?, updated_at = NOW()",
		haveread.UserID, haveread.ChannelID, haveread.MessageID, haveread.MessageID)
	if err != nil {
		return err
	}
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

	if err != nil {
		fmt.Println("ERROR: ", err)
		return 0, nil
	}
	return haveread.MessageID, nil
}

type MessageCount struct {
	Count     int64 `db:"cnt"`
	ChannelID int64 `db:"channel_id"`
}

func FetchUnreadMessageCount(userID int64) (map[string]MessageCount, error) {
	var mc []MessageCount
	err := db.Select(&mc, "SELECT COUNT(*) as cnt, channel_id FROM message GROUP BY channel_id")
	if err != nil {
		return nil, err
	}
	fmt.Println(mc)
	resp := map[string]MessageCount{}

	for _, v := range mc {
		resp[strconv.Itoa(int(v.ChannelID))] = v
	}

	var mc2 []MessageCount
	err = db.Select(&mc2, "SELECT COUNT(*) as cnt, m.channel_id FROM message m INNER JOIN haveread h ON h.channel_id = m.channel_id AND h.message_id <= m.id AND h.user_id = ? GROUP BY m.channel_id", userID)
	if err != nil {
		return nil, err
	}
	for _, v := range mc2 {
		v.Count--
		resp[strconv.Itoa(int(v.ChannelID))] = v
	}
	fmt.Println(mc2)

	return resp, nil
}
