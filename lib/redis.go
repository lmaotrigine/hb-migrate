package lib

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
)

type RedisClient struct {
	handler *rejson.Handler
	client  *redis.Client
}

func NewRedisClient(addr string, pass string) (*RedisClient, error) {
	rh := rejson.NewReJSONHandler()
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
	})
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}
	rh.SetGoRedisClient(client)
	return &RedisClient{handler: rh, client: client}, nil
}

func (c *RedisClient) GetAllBeats() ([]BeatLegacy, error) {
	res, err := c.handler.JSONGet("beats", ".")
	if err != nil {
		return nil, err
	}
	var beats []BeatLegacy
	err = json.Unmarshal([]byte(res.([]uint8)), &beats)
	if err != nil {
		return nil, err
	}
	return beats, nil
}

func (c *RedisClient) GetAllDevices() ([]DeviceLegacy, error) {
	res, err := c.handler.JSONGet("devices", ".")
	if err != nil {
		return nil, err
	}
	var devices []DeviceLegacy
	err = json.Unmarshal([]byte(res.([]uint8)), &devices)
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func (c *RedisClient) GetStats() (*StatsLegacy, error) {
	res, err := c.handler.JSONGet("stats", ".")
	if err != nil {
		return nil, err
	}
	var stats StatsLegacy
	err = json.Unmarshal([]byte(res.([]uint8)), &stats)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}
