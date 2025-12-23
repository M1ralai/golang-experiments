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
	ctx, cancel := context.WithCancel(context.Background())

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})

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
				Count:    1,
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
				payload := entry.Values["event_data"].(string)

				if err := handler([]byte(payload)); err == nil {
					r.client.XAck(ctx, topic, r.group, entry.ID)
				} else {
					log.Printf("Handler error for message %s: %v", entry.ID, err)
				}
			}
		}
	}
}

func (r *redisBus) processPendingMessages(ctx context.Context, topic string, handler func([]byte) error) {
	log.Printf("Processing pending messages for topic: %s", topic)

	pending, err := r.client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    r.group,
		Consumer: r.worker,
		Streams:  []string{topic, "0"},
		Count:    100,
	}).Result()

	if err != nil && err != redis.Nil {
		log.Printf("Failed to read pending messages: %v", err)
		return
	}

	if len(pending) > 0 {
		for _, entry := range pending[0].Messages {
			payload := entry.Values["event_data"].(string)

			if err := handler([]byte(payload)); err == nil {
				r.client.XAck(ctx, topic, r.group, entry.ID)
				log.Printf("Processed pending message: %s", entry.ID)
			} else {
				log.Printf("Failed to process pending message %s: %v", entry.ID, err)
			}
		}
	}
}

func (r *redisBus) Close() error {
	r.cancel()
	return r.client.Close()
}
