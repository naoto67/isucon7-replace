package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type ChannelInfo struct {
	ID          int64     `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	UpdatedAt   time.Time `db:"updated_at"`
	CreatedAt   time.Time `db:"created_at"`
}

const (
	CHANNEL_INDEX = "channels"
)

func queryChannels() ([]int64, error) {
	res := []int64{}
	channels, err := fetchChannels()
	if err != nil {
		fmt.Println("DEBUG: use mysql data in queryChannels", err)
		err = db.Select(&res, "SELECT id FROM channel")
	} else {
		for _, v := range channels {
			res = append(res, v.ID)
		}
	}
	return res, err
}

func initChannels() error {
	var channels []ChannelInfo
	err := db.Select(&channels, "SELECT * FROM channel")
	if err != nil {
		return err
	}
	for _, v := range channels {
		AddChannelCache(v)
	}

	return err
}

func fetchChannels() ([]ChannelInfo, error) {
	data, err := cache.LFetchAll(CHANNEL_INDEX)
	if err != nil {
		return nil, err
	}
	var channels []ChannelInfo
	for i := range data {
		var ch ChannelInfo
		err = json.Unmarshal(data[i], &ch)
		if err != nil {
			fmt.Println("WARN: ", err)
			return nil, err
		}
		channels = append(channels, ch)
	}
	return channels, nil
}

func AddChannelCache(ch ChannelInfo) error {
	err := cache.LPush(CHANNEL_INDEX, ch)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s:%s", MESSAGE_COUNT_PREFIX, ch.ID)
	return cache.Set(key, 0)
}
