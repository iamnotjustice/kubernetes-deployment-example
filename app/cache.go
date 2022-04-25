package main

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

func connectToCache() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	pong, err := rdb.Ping(context.Background()).Result()
	if err == nil {
		logrus.Info(pong, err)
	} else {
		logrus.Error(err)
	}
}

func getFromCache(ipaddress string) (*IPGeoData, error) {
	cacheData, err := rdb.Get(context.Background(), ipaddress).Result()
	if err != nil {
		return nil, err
	}

	var obj IPGeoData
	if err := json.Unmarshal([]byte(cacheData), &obj); err != nil {
		return nil, err
	}

	return &obj, nil
}

func saveToCache(data IPGeoData) {
	//store in redis
	err := rdb.Set(context.Background(), data.IP, data, time.Hour*24).Err()
	if err != nil {
		logrus.Error(err)
	}
}
