package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

type Redis struct {
	Client  *redis.Client
	Context context.Context
}

func newRedis() (Redis, error) {
	redisAddress := os.Getenv("REDIS_ADDRESS")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDatabase, redisDatabaseError := strconv.Atoi(os.Getenv("REDIS_DATABASE"))
	redisPort := os.Getenv("REDIS_PORT")
	if redisDatabaseError.Error() != "" {
		log.Info("could not find REDIS_DATABASE so using default value 0")
		redisDatabase = 0
	}
	if redisPort == "" {
		redisPort = "6379"
	}
	if redisAddress == "" || redisPassword == "" {
		return Redis{}, fmt.Errorf("REDIS_ADDRESS or REDIS_PASSWORD is empty, cannot initialize")
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisAddress, redisPort),
		Password: "",
		DB:       redisDatabase,
	})
	return Redis{Client: rdb, Context: context.Background()}, nil
}

func (r *Redis) getTeams() ([]string, error) {
	teams, _ := r.Client.Keys(r.Context, "*").Result()
	for _, t := range teams {
		fmt.Println(t)
	}
	var cursor uint64
	keys, cursor, err := r.Client.Scan(r.Context, cursor, "prefix:*", 0).Result()
	if err != nil {
		log.Error(err)
	}
	for _, k := range keys {
		fmt.Println(k)
	}
	return nil, nil
}

func (r *Redis) getKeys(key string) ([]string, error) {
	keys, err := r.Client.Keys(r.Context, fmt.Sprintf("%s*", key)).Result()
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func (r *Redis) getValue(key string) (string, error) {
	value, err := r.Client.Get(r.Context, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

func (r *Redis) addKey(key string, value interface{}) error {
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(value)
	err := r.Client.Set(r.Context, key, reqBodyBytes.Bytes(), 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) removeKey(key string) error {
	keys, err := r.getKeys(key)
	if err != nil {
		return err
	}
	pipe := r.Client.Pipeline()
	for _, key := range keys {
		pipe.Del(r.Context, key)
	}
	pipe.Exec(r.Context)
	return nil
}
