package eventbus

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	MaxRetries = 3
	DLQSuffix  = "_dlq"
	BatchSize  = 50
)

type EventBus interface {
	Publish(ctx context.Context, topic string, payload any) error
	Subscribe(ctx context.Context, topic string, handler func(payload []byte) error)
	Close() error
}

type redisBus struct {
	client *redis.Client
	group  string
	worker string
	ctx    context.Context
	cancel context.CancelFunc
}

func NewRedisBus(addr, password, groupName string) EventBus {
	if addr == "" {
		log.Fatal("REDIS_ADDR is required")
	}

	ctx, cancel := context.WithCancel(context.Background())

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "unknown"
	}
	workerName := fmt.Sprintf("%s-%d", hostname, os.Getpid())

	return &redisBus{
		client: rdb,
		group:  groupName,
		worker: workerName,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (r *redisBus) Publish(ctx context.Context, topic string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	return r.client.XAdd(ctx, &redis.XAddArgs{
		Stream: topic,
		Values: map[string]any{"event_data": data},
	}).Err()
}

func (r *redisBus) Subscribe(ctx context.Context, topic string, handler func(payload []byte) error) {
	r.client.XGroupCreateMkStream(ctx, topic, r.group, "0").Err()

	go r.listenLoop(ctx, topic, handler)
}

func (r *redisBus) listenLoop(ctx context.Context, topic string, handler func([]byte) error) {
	log.Printf("Redis Stream listening: Topic=%s Group=%s Worker=%s", topic, r.group, r.worker)

	r.processPendingMessages(ctx, topic, handler)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Stopping listener for topic: %s", topic)
			return
		default:
			entries, err := r.client.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    r.group,
				Consumer: r.worker,
				Streams:  []string{topic, ">"},
				Count:    BatchSize,
				Block:    time.Second,
			}).Result()

			if err != nil {
				if err == redis.Nil || err == context.Canceled {
					continue
				}
				log.Printf("XReadGroup error: %v", err)
				time.Sleep(time.Second)
				continue
			}

			for _, entry := range entries[0].Messages {
				r.handleMessage(ctx, topic, entry, handler)
			}
		}
	}
}

func (r *redisBus) handleMessage(ctx context.Context, topic string, entry redis.XMessage, handler func([]byte) error) {
	payload := entry.Values["event_data"].(string)

	err := handler([]byte(payload))
	if err == nil {
		r.client.XAck(ctx, topic, r.group, entry.ID)
		return
	}

	log.Printf("Handler error for message %s: %v", entry.ID, err)

	pending, _ := r.client.XPending(ctx, topic, r.group).Result()
	if pending != nil {
		for i := 0; i < MaxRetries; i++ {
			time.Sleep(time.Duration(i+1) * time.Second)
			if err := handler([]byte(payload)); err == nil {
				r.client.XAck(ctx, topic, r.group, entry.ID)
				log.Printf("✓ Retry %d succeeded for message %s", i+1, entry.ID)
				return
			}
		}
	}

	r.moveToDLQ(ctx, topic, entry)
	r.client.XAck(ctx, topic, r.group, entry.ID)
}

func (r *redisBus) moveToDLQ(ctx context.Context, topic string, entry redis.XMessage) {
	dlqTopic := topic + DLQSuffix

	err := r.client.XAdd(ctx, &redis.XAddArgs{
		Stream: dlqTopic,
		Values: entry.Values,
	}).Err()

	if err != nil {
		log.Printf("Failed to move message %s to DLQ: %v", entry.ID, err)
	} else {
		log.Printf("⚠️  Moved message %s to DLQ: %s", entry.ID, dlqTopic)
	}
}

func (r *redisBus) processPendingMessages(ctx context.Context, topic string, handler func([]byte) error) {
	log.Printf("Processing pending messages for topic: %s", topic)

	start := "0"
	for {
		pending, err := r.client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    r.group,
			Consumer: r.worker,
			Streams:  []string{topic, start},
			Count:    BatchSize,
		}).Result()

		if err != nil && err != redis.Nil {
			log.Printf("Failed to read pending messages: %v", err)
			return
		}

		if len(pending) == 0 || len(pending[0].Messages) == 0 {
			break
		}

		for _, entry := range pending[0].Messages {
			r.handleMessage(ctx, topic, entry, handler)
			start = entry.ID
		}

		if len(pending[0].Messages) < BatchSize {
			break
		}
	}

	log.Printf("✓ Pending messages processed for topic: %s", topic)
}

func (r *redisBus) Close() error {
	r.cancel()
	return r.client.Close()
}
