package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gomodule/redigo/redis"
)

type Cache struct {
	MemcacheClient *memcache.Client
	RedisPool      *redis.Pool
}

var cache *Cache

func NewCache(redisServer, memcacheServer string) {
	cache.RedisPool = &redis.Pool{
		MaxIdle:     6,
		MaxActive:   3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", redisServer) },
	}
	conn := cache.RedisPool.Get()
	_, err := conn.Do("PING")
	if err != nil {
		panic(err)
	}

	cache.MemcacheClient = memcache.New(memcacheServer)
	err = cache.MemcacheClient.Ping()
	if err != nil {
		fmt.Println("failed ping.", err)
		panic(err)
	}
	fmt.Println("PONG!!!")

	return
}

func (cache *Cache) Set(key string, value interface{}) error {
	data, _ := json.Marshal(value)
	return cache.MemcacheClient.Set(&memcache.Item{Key: key, Value: data})
}

func (cache *Cache) Get(key string) ([]byte, error) {
	item, err := cache.MemcacheClient.Get(key)
	if err != nil {
		return nil, err
	}
	return item.Value, nil
}

func (cache *Cache) FlushAll() {
	cache.MemcacheClient.FlushAll()
	conn := cache.RedisPool.Get()
	conn.Do("FLUSHALL")
}

func (cache *Cache) LPush(key string, value interface{}) error {
	data, _ := json.Marshal(value)
	conn := cache.RedisPool.Get()
	_, err := conn.Do("LPUSH", key, data)
	return err
}

func (cache *Cache) LFetchAll(key string) ([][]byte, error) {
	conn := cache.RedisPool.Get()
	data, err := redis.ByteSlices(conn.Do("LRANGE", key, 0, -1))
	return data, err
}
