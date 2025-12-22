# Task Module API Documentation

## POST /api/tasks
Yeni bir task oluşturur.

### Request Body
```json
{
  "title": "string" // Zorunlu (1-255 karakter)
}
```

### Response Body (Success - 201)
```json
{
  "success": true,
  "message": "Task başarıyla oluşturuldu",
  "data": {
    "id": "uuid",
    "title": "string",
    "status": "todo",
    "created_by": "uuid",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  },
  "error": null,
  "timestamp": "string"
}
```

### Validation Rules
- **title**: Zorunlu (required), min 1, max 255 karakter

---

## GET /api/tasks
Tüm task'ları listeler.

### Response Body (Success - 200)
```json
{
  "success": true,
  "message": "Task listesi başarıyla getirildi",
  "data": [
    {
      "id": "uuid",
      "title": "string",
      "status": "todo|in_progress|done",
      "created_by": "uuid",
      "created_at": "timestamp",
      "updated_at": "timestamp"
    }
  ],
  "error": null,
  "timestamp": "string"
}
```

---

## GET /api/tasks/{id}
ID'ye göre task detayını getirir.

### Response Body (Success - 200)
```json
{
  "success": true,
  "message": "Task başarıyla getirildi",
  "data": {
    "id": "uuid",
    "title": "string",
    "status": "todo|in_progress|done",
    "created_by": "uuid",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  },
  "error": null,
  "timestamp": "string"
}
```

### Response Body (Error - 404)
```json
{
  "success": false,
  "message": "Task bulunamadı",
  "data": null,
  "error": {
    "code": "NOT_FOUND",
    "details": ""
  },
  "timestamp": "string"
}
```

---

## PATCH /api/tasks/{id}/status
Task durumunu günceller.

### Request Body
```json
{
  "status": "string" // Zorunlu: "todo", "in_progress", veya "done"
}
```

### Response Body (Success - 200)
```json
{
  "success": true,
  "message": "Task durumu başarıyla güncellendi",
  "data": null,
  "error": null,
  "timestamp": "string"
}
```

### Validation Rules
- **status**: Zorunlu (required), değerler: `todo`, `in_progress`, `done`

---

## POST /api/tasks/{id}/assignments
Task'a kullanıcı atar.

### Request Body
```json
{
  "user_id": "uuid" // Zorunlu
}
```

### Response Body (Success - 201)
```json
{
  "success": true,
  "message": "Task başarıyla atandı",
  "data": {
    "id": "uuid",
    "task_id": "uuid",
    "user_id": "uuid",
    "created_at": "timestamp"
  },
  "error": null,
  "timestamp": "string"
}
```

### Validation Rules
- **user_id**: Zorunlu (required), geçerli UUID formatında

---

## GET /api/tasks/{id}/assignments
Task'a atanmış kullanıcıları listeler.

### Response Body (Success - 200)
```json
{
  "success": true,
  "message": "Task atamaları başarıyla getirildi",
  "data": [
    {
      "id": "uuid",
      "task_id": "uuid",
      "user_id": "uuid",
      "created_at": "timestamp"
    }
  ],
  "error": null,
  "timestamp": "string"
}
```

---

## DELETE /api/tasks/assignments/{id}
Task atamasını kaldırır.

### Response Body (Success - 200)
```json
{
  "success": true,
  "message": "Task ataması başarıyla kaldırıldı",
  "data": null,
  "error": null,
  "timestamp": "string"
}
```
