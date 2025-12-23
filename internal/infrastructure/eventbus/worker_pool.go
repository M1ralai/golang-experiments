package eventbus

import (
	"context"
	"log"
	"sync"
)

type EventPublishJob struct {
	Topic   string
	Payload any
}

type WorkerPool struct {
	bus      EventBus
	jobQueue chan EventPublishJob
	workers  int
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewWorkerPool(bus EventBus, workers int, queueSize int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	pool := &WorkerPool{
		bus:      bus,
		jobQueue: make(chan EventPublishJob, queueSize),
		workers:  workers,
		ctx:      ctx,
		cancel:   cancel,
	}

	pool.start()
	return pool
}

func (p *WorkerPool) start() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
}

func (p *WorkerPool) worker(id int) {
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			log.Printf("Worker %d shutting down", id)
			return
		case job := <-p.jobQueue:
			if err := p.bus.Publish(p.ctx, job.Topic, job.Payload); err != nil {
				log.Printf("Worker %d: Failed to publish event to %s: %v", id, job.Topic, err)
			}
		}
	}
}

func (p *WorkerPool) PublishAsync(topic string, payload any) error {
	select {
	case p.jobQueue <- EventPublishJob{Topic: topic, Payload: payload}:
		return nil
	default:
		return ErrQueueFull
	}
}

func (p *WorkerPool) Shutdown() {
	log.Println("Shutting down worker pool...")
	p.cancel()
	p.wg.Wait()
	close(p.jobQueue)
	log.Println("Worker pool shut down complete")
}

var ErrQueueFull = &QueueFullError{}

type QueueFullError struct{}

func (e *QueueFullError) Error() string {
	return "event queue is full"
}
