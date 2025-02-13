package redis

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/Mark-Grigorev/auth/internal/model"
	"github.com/go-redis/redis/v8"
)

type Client struct {
	ctx context.Context
	rdb redis.UniversalClient
	ttl time.Duration
}

func New(cfg model.RedisConfig) (*Client, error) {
	ctx := context.Background()
	addrs := strings.Split(cfg.Servers, ",")
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    addrs,
		Password: cfg.Password,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		return &Client{}, err
	}
	return &Client{
		ctx: ctx,
		rdb: rdb,
		ttl: time.Duration(cfg.TTL) * time.Minute,
	}, nil
}

func (c *Client) SaveToken(userID int64, token string) error {
	pipe := c.rdb.Pipeline()
	key := c.userIDKey(userID)

	_ = pipe.HSet(c.ctx, key, "user_id", userID)
	_ = pipe.HSet(c.ctx, key, "token", token)
	_ = pipe.Expire(c.ctx, key, c.ttl)
	_, err := pipe.Exec(c.ctx)
	_ = pipe.Close()

	return err
}

func (c *Client) GetTokenByUserID(userID int64) (string, error) {
	key := c.userIDKey(userID)
	val, err := c.rdb.HGetAll(c.ctx, key).Result()
	if err != nil {
		return "", err
	}
	if len(val) == 0 {
		return "", errors.New("token not found in redis")
	}
	token := val["token"]
	return token, nil
}

func (c *Client) userIDKey(userID int64) string {
	return "token_" + strconv.FormatInt(userID, 10)
}
