package eventbus

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type EventBus interface {
	Publish(ctx context.Context, topic string, payload any) error
	Subscribe(topic string, handler func(payload []byte) error)
	Close() error
}

type redisBus struct {
	client *redis.Client
	group  string
	worker string
}

func NewRedisBus(addr, password, groupName string) EventBus {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})

	workerName := fmt.Sprintf("consumer-%d", time.Now().UnixNano())

	return &redisBus{
		client: rdb,
		group:  groupName,
		worker: workerName,
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

func (r *redisBus) Subscribe(topic string, handler func(payload []byte) error) {
	// Consumer Group oluştur (Hata verirse yoksay, çünkü zaten varsa hata verir)
	r.client.XGroupCreateMkStream(context.Background(), topic, r.group, "0").Err()

	go r.listenLoop(topic, handler)
}

func (r *redisBus) listenLoop(topic string, handler func([]byte) error) {
	log.Printf("Redis Stream Dinleniyor: Topic=%s Group=%s", topic, r.group)

	for {
		entries, err := r.client.XReadGroup(context.Background(), &redis.XReadGroupArgs{
			Group:    r.group,
			Consumer: r.worker,
			Streams:  []string{topic, ">"},
			Count:    1,
			Block:    0,
		}).Result()

		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		for _, entry := range entries[0].Messages {
			payload := entry.Values["event_data"].(string)

			if err := handler([]byte(payload)); err == nil {
				r.client.XAck(context.Background(), topic, r.group, entry.ID)
			} else {
				log.Printf("işlem hatası: %v", err)
			}
		}
	}
}

func (r *redisBus) Close() error {
	return r.client.Close()
}
