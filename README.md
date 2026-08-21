# Go Backend Architecture Experiments

A modular Go backend used to explore event-driven application patterns around a task-assignment domain. It combines a conventional HTTP/service/repository stack with PostgreSQL, Redis Streams, a transactional outbox, and asynchronous consumers.

## What this project demonstrates

- Domain-oriented modules for authentication, users, tasks, and notifications
- PostgreSQL repositories and embedded SQL migrations
- A transactional outbox written in the same database transaction as task changes
- A background processor that publishes pending outbox records to Redis Streams
- Consumer-group-based event handling and a notification listener
- A bounded in-process worker pool for asynchronous publish jobs
- JWT middleware, Prometheus HTTP metrics, and Zap logging

## Architecture

```text
HTTP request
  -> handler -> task service -> PostgreSQL transaction
                              -> task/assignment changes
                              -> outbox event

outbox processor -> Redis Stream -> consumer group -> notification listener
```

The HTTP path keeps module boundaries similar to a modular monolith: handlers call services, services use repository contracts, and PostgreSQL implementations own SQL. For task-assignment events, the service stores the domain change and outbox record together. A background processor polls unprocessed rows and publishes them to Redis; consumers acknowledge successful work and failed messages can be moved to a dead-letter stream.

The notification listener currently demonstrates the consumer boundary by logging an email-shaped notification rather than integrating a mail provider.

## Run locally

Requirements: Go 1.25+, PostgreSQL, and Redis.

```bash
git clone https://github.com/M1ralai/golang-experiments.git
cd golang-experiments
cp .env.example .env
go mod download
go run cmd/api/main.go
```

Set the PostgreSQL connection, `JWT_SECRET`, and `REDIS_ADDR` in `.env`. Migrations run during startup.

After startup, infrastructure endpoints include:

```bash
curl http://localhost:8080/health
curl http://localhost:8080/metrics
```

## Limitations

- This repository is an architecture experiment, not a production-ready service.
- The notification consumer logs a simulated email; it does not send one.
- Outbox retries, pending-message recovery, and dead-letter behavior have not been validated under sustained failure or load.
- Automated tests and deployment documentation are limited.
- The module path still uses the original template name, so extracting the experiment as a reusable package would require cleanup.
