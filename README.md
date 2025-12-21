# Go Modüler Monolit Template

JWT kimlik doğrulama, PostgreSQL veritabanı, yapısal loglama ve metrikler ile üretime hazır, temiz bir Go API şablonu.

## 🏗️ Mimari

```
├── cmd/api/                    # Uygulama giriş noktası
│   └── main.go                 # Bootstrap & lifecycle yönetimi
├── internal/
│   ├── app/
│   │   └── server.go           # HTTP sunucu & routing
│   ├── common/
│   │   ├── stype/              # Paylaşılan tipler (API response formatı)
│   │   ├── utils/              # Yardımcı fonksiyonlar (JSON, response writers)
│   │   └── validation/         # Request validasyonu (go-playground)
│   ├── infrastructure/
│   │   ├── database/           # PostgreSQL bağlantısı & migration'lar
│   │   ├── logger/             # Zap yapısal loglama (DB'ye kayıt)
│   │   ├── metrics/            # Prometheus metrikleri
│   │   └── middleware/         # Auth, recovery, timeout, metrics middleware
│   └── modules/
│       ├── auth/               # JWT kimlik doğrulama (login)
│       ├── health/             # Health check endpoint
│       └── user/               # Kullanıcı CRUD işlemleri
└── go.mod
```

## 🚀 Özellikler

- **JWT Kimlik Doğrulama** - Rol destekli güvenli token tabanlı auth
- **Request Validasyonu** - go-playground/validator ile Türkçe çeviriler
- **Veritabanı Migration'ları** - golang-migrate ile başlangıçta otomatik migration
- **Yapısal Loglama** - Veritabanına kayıt yapan Zap logger
- **Prometheus Metrikleri** - `/metrics` endpoint'inde hazır metrikler
- **Graceful Shutdown** - Düzgün sinyal yönetimi ve temizlik
- **Middleware Yığını** - Recovery, timeout, auth ve metrics middleware
- **Temiz Mimari** - Domain → Repository → Service → HTTP katmanları

## 📋 Gereksinimler

- Go 1.21+
- PostgreSQL 14+

## 🛠️ Kurulum

1. Repository'yi klonla
2. Ortam dosyasını kopyala:
   ```bash
   cp .env.example .env
   ```
3. `.env` dosyasını yapılandır:
   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=myapp
   JWT_SECRET=cok-gizli-anahtar-bunu-degistir
   API_PORT=8080
   ```
4. Uygulamayı çalıştır:
   ```bash
   go run cmd/api/main.go
   ```

## 📡 API Endpoint'leri

### Public Route'lar

| Metod | Endpoint  | Açıklama              |
|-------|-----------|----------------------|
| POST  | /login    | Kullanıcı girişi     |
| GET   | /health   | Sağlık kontrolü      |
| GET   | /metrics  | Prometheus metrikleri|

### Korumalı Route'lar (JWT Gerekli)

| Metod  | Endpoint        | Açıklama           |
|--------|-----------------|-------------------|
| GET    | /api/users      | Tüm kullanıcıları listele |
| POST   | /api/users      | Yeni kullanıcı oluştur   |
| DELETE | /api/users/{id} | Kullanıcı sil           |

## 🔧 Yeni Modül Ekleme

Katmanlı yapıyı takip et:

1. **Domain** (`internal/modules/moduladi/domain/`)
   - `entity.go` - Veri yapıları
   - `repository.go` - Repository interface'i

2. **Repository** (`internal/modules/moduladi/repository/`)
   - `pg_repository.go` - PostgreSQL implementasyonu

3. **Service** (`internal/modules/moduladi/service/`)
   - `service.go` - İş mantığı

4. **HTTP** (`internal/modules/moduladi/http/`)
   - `handler.go` - HTTP handler'ları

5. `internal/app/server.go` dosyasında bağla

## 📦 Teknoloji Yığını

- **Router**: gorilla/mux
- **Veritabanı**: sqlx + lib/pq
- **Migration**: golang-migrate
- **Auth**: golang-jwt
- **Validasyon**: go-playground/validator
- **Loglama**: uber/zap
- **Metrikler**: prometheus/client_golang
- **Şifreleme**: bcrypt

## 📄 Lisans

MIT
