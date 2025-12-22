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
│       ├── task/               # Task yönetimi (CRUD + atama)
│       └── user/               # Kullanıcı CRUD işlemleri
└── go.mod
```

## 🚀 Özellikler

- **JWT Kimlik Doğrulama** - Rol ve user_id destekli güvenli token tabanlı auth
- **UUID Primary Keys** - Tüm tablolarda UUID kullanımı
- **Request Validasyonu** - go-playground/validator ile Türkçe çeviriler
- **Veritabanı Migration'ları** - golang-migrate ile başlangıçta otomatik migration
- **Yapısal Loglama** - Veritabanına kayıt yapan Zap logger
- **Prometheus Metrikleri** - `/metrics` endpoint'inde hazır metrikler
- **Graceful Shutdown** - Düzgün sinyal yönetimi ve temizlik
- **Middleware Yığını** - Recovery, timeout, auth ve metrics middleware
- **Temiz Mimari** - Domain → Repository → Service → HTTP katmanları
- **Task Modülü** - Task yönetimi, kullanıcı ataması ve aktivite takibi

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

#### User Modülü

| Metod  | Endpoint        | Açıklama                |
|--------|-----------------|------------------------|
| GET    | /api/users      | Tüm kullanıcıları listele |
| POST   | /api/users      | Yeni kullanıcı oluştur   |
| DELETE | /api/users/{id} | Kullanıcı sil           |

#### Task Modülü

| Metod  | Endpoint                       | Açıklama                    |
|--------|--------------------------------|----------------------------|
| GET    | /api/tasks                     | Tüm task'ları listele       |
| POST   | /api/tasks                     | Yeni task oluştur          |
| GET    | /api/tasks/{id}                | Task detayını getir        |
| PATCH  | /api/tasks/{id}/status         | Task durumunu güncelle     |
| GET    | /api/tasks/{id}/assignments    | Task atamalarını listele   |
| POST   | /api/tasks/{id}/assignments    | Task'a kullanıcı ata       |
| DELETE | /api/tasks/assignments/{id}    | Task atamasını kaldır      |

## 🔧 Yeni Modül Ekleme

Katmanlı yapıyı takip et:

1. **Domain** (`internal/modules/moduladi/domain/`)
   - `entity.go` - Veri yapıları (JSON/DB tag'leri ile)
   - `repository.go` - Repository interface'i

2. **Repository** (`internal/modules/moduladi/repository/`)
   - `pg_repository.go` - PostgreSQL implementasyonu

3. **Service** (`internal/modules/moduladi/service/`)
   - `service.go` - İş mantığı (infrastructure logger ile)

4. **HTTP** (`internal/modules/moduladi/http/`)
   - `handler.go` - HTTP handler'ları

5. **Migration** (`internal/infrastructure/database/migrations/`)
   - `000XXX_create_xxx_tables.up.sql` - Tablo oluşturma
   - `000XXX_create_xxx_tables.down.sql` - Rollback

6. **Entegrasyon**
   - `internal/app/server.go` dosyasında repo, service ve handler'ı bağla
   - Route'ları ekle

7. **Dokümantasyon**
   - `api.md` - Endpoint dokümantasyonu

## 📦 Teknoloji Yığını

- **Router**: gorilla/mux
- **Veritabanı**: sqlx + lib/pq
- **Migration**: golang-migrate
- **Auth**: golang-jwt
- **Validasyon**: go-playground/validator
- **Loglama**: uber/zap
- **Metrikler**: prometheus/client_golang
- **Şifreleme**: bcrypt

## 📁 Modül Yapısı

Her modül aşağıdaki yapıyı takip eder:

```
modules/
└── moduladi/
    ├── api.md              # API dokümantasyonu
    ├── domain/
    │   ├── entity.go       # Domain entity'leri
    │   └── repository.go   # Repository interface'leri
    ├── repository/
    │   └── pg_repository.go # PostgreSQL implementasyonu
    ├── service/
    │   └── service.go      # İş mantığı katmanı
    └── http/
        └── handler.go      # HTTP handler'ları
```

## 📄 Lisans

MIT
