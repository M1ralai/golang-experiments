package outbox

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/eventbus"
)

type Processor struct {
	repo      Repository
	eventBus  eventbus.EventBus
	interval  time.Duration
	batchSize int
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewProcessor(repo Repository, eventBus eventbus.EventBus, interval time.Duration, batchSize int) *Processor {
	ctx, cancel := context.WithCancel(context.Background())
	return &Processor{
		repo:      repo,
		eventBus:  eventBus,
		interval:  interval,
		batchSize: batchSize,
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (p *Processor) Start() {
	log.Println("✓ Outbox processor started")

	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-p.ctx.Done():
			log.Println("✓ Outbox processor stopped")
			return
		case <-ticker.C:
			p.processEvents()
		}
	}
}

func (p *Processor) Stop() {
	p.cancel()
}

func (p *Processor) processEvents() {
	events, err := p.repo.GetUnprocessed(p.ctx, p.batchSize)
	if err != nil {
		log.Printf("Failed to fetch outbox events: %v", err)
		return
	}

	if len(events) == 0 {
		return
	}

	log.Printf("Processing %d outbox events", len(events))

	for _, event := range events {
		if err := p.publishEvent(event); err != nil {
			log.Printf("Failed to publish event %s: %v", event.ID, err)
			p.repo.MarkFailed(p.ctx, event.ID, err.Error())
		} else {
			p.repo.MarkProcessed(p.ctx, event.ID)
			log.Printf("✓ Published event %s (type: %s)", event.ID, event.EventType)
		}
	}
}

func (p *Processor) publishEvent(event *OutboxEvent) error {
	var payload map[string]interface{}
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	topic := event.EventType
	return p.eventBus.Publish(p.ctx, topic, payload)
}
