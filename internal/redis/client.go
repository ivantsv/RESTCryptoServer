package redis

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/redis/go-redis/v9"
)

type PriceHistoryEntry struct {
    Price     float64   `json:"price"`
    Timestamp time.Time `json:"timestamp"`
}

type RedisClient struct {
    client *redis.Client
    ctx    context.Context
}

func NewRedisClient() (*RedisClient, error) {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("REDIS_PORT")
    if port == "" {
        port = "6379"
    }

	password := os.Getenv("REDIS_PASSWORD")

	rdb := redis.NewClient(&redis.Options{
        Addr:     host + ":" + port,
        Password: password,
        DB:       0,
    })

	ctx := context.Background()

	_, err := rdb.Ping(ctx).Result()
    if err != nil {
        return nil, fmt.Errorf("failed to connect to Redis: %w", err)
    }

	log.Println("Successfully connected to Redis")
    
    return &RedisClient{
        client: rdb,
        ctx:    ctx,
    }, nil
}

func (r *RedisClient) AddPriceHistory(symbol string, price float64) error {
	key := fmt.Sprintf("price_history:%s", symbol)

	entry := PriceHistoryEntry{
        Price:     price,
        Timestamp: time.Now().UTC(),
    }

	data, err := json.Marshal(entry)
    if err != nil {
        return fmt.Errorf("failed to marshal price entry: %w", err)
    }

	err = r.client.LPush(r.ctx, key, string(data)).Err()
    if err != nil {
        return fmt.Errorf("failed to add price to history: %w", err)
    }

	err = r.client.LTrim(r.ctx, key, 0, 99).Err()
    if err != nil {
        return fmt.Errorf("failed to trim price history: %w", err)
    }
    
    log.Printf("Added price history for %s: $%.2f", symbol, price)
    return nil
}

func (r *RedisClient) GetLatestPrice(symbol string) (*PriceHistoryEntry, error) {
    key := fmt.Sprintf("price_history:%s", symbol)
    
    result, err := r.client.LIndex(r.ctx, key, 0).Result()
    if err != nil {
        if err == redis.Nil {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to get latest price: %w", err)
    }
    
    var entry PriceHistoryEntry
    if err := json.Unmarshal([]byte(result), &entry); err != nil {
        return nil, fmt.Errorf("failed to unmarshal latest price: %w", err)
    }
    
    return &entry, nil
}

func (r *RedisClient) DeletePriceHistory(symbol string) error {
    key := fmt.Sprintf("price_history:%s", symbol)
    
    err := r.client.Del(r.ctx, key).Err()
    if err != nil {
        return fmt.Errorf("failed to delete price history: %w", err)
    }
    
    log.Printf("Deleted price history for %s", symbol)
    return nil
}

func (r *RedisClient) GetHistoryCount(symbol string) (int, error) {
    key := fmt.Sprintf("price_history:%s", symbol)
    
    count, err := r.client.LLen(r.ctx, key).Result()
    if err != nil {
        return 0, fmt.Errorf("failed to get history count: %w", err)
    }
    
    return int(count), nil
}

func (r *RedisClient) GetPriceHistory(symbol string, limit int) ([]PriceHistoryEntry, error) {
    if limit <= 0 || limit > 100 {
        limit = 100
    }
    
    key := fmt.Sprintf("price_history:%s", symbol)

    result, err := r.client.LRange(r.ctx, key, 0, int64(limit-1)).Result()
    if err != nil {
        return nil, fmt.Errorf("failed to get price history: %w", err)
    }
    
    var history []PriceHistoryEntry
    for _, data := range result {
        var entry PriceHistoryEntry
        if err := json.Unmarshal([]byte(data), &entry); err != nil {
            log.Printf("Failed to unmarshal price entry: %v", err)
            continue
        }
        history = append(history, entry)
    }
    
    return history, nil
}

func (r *RedisClient) Close() error {
    if r.client != nil {
        return r.client.Close()
    }
    return nil
}