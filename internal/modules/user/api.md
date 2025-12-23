# User Module API Documentation

## GET /api/users
Tüm kayıtlı kullanıcıları listeler.

### Response Body (Success - 200)
```json
{
  "success": true,
  "message": "Kullanıcılar başarıyla getirildi",
  "data": [
    {
      "id": "uuid",
      "username": "string",
      "role": "string",
      "ad": "string",
      "soyad": "string",
      "telefon": "string",
      "email": "string"
    }
  ],
  "error": null,
  "timestamp": "string"
}
```

---

## POST /api/users
Yeni bir sistem kullanıcısı oluşturur.

### Request Body
```json
{
  "username": "string", // Zorunlu
  "password": "string", // Zorunlu (min 3 karakter)
  "role": "string",     // Zorunlu (ADMIN, SECRETARY vb.)
  "ad": "string",       // Opsiyonel
  "soyad": "string",    // Opsiyonel
  "telefon": "string",  // Opsiyonel
  "email": "string"     // Opsiyonel
}
```

### Response Body (Success - 201)
```json
{
  "success": true,
  "message": "Kullanıcı başarıyla eklendi",
  "data": {
    "id": "uuid",
    "username": "string",
    "role": "string",
    "ad": "string",
    "soyad": "string",
    "telefon": "string",
    "email": "string"
  },
  "error": null,
  "timestamp": "string"
}
```

### Validation Rules
- **username**: Zorunlu (required)
- **password**: Zorunlu (required) ve en az 3 karakter
- **role**: Zorunlu (required)

---

## GET /api/users/{id}
UUID'ye göre tek bir kullanıcıyı getirir. Kullanıcının username ve email bilgilerini içerir.

### Path Parameters
- **id**: Kullanıcı UUID'si

### Response Body (Success - 200)
```json
{
  "success": true,
  "message": "Kullanıcı başarıyla getirildi",
  "data": {
    "id": "uuid",
    "username": "string",
    "ad": "string",
    "soyad": "string",
    "email": "string"
  },
  "error": null,
  "timestamp": "string"
}
```

### Response Body (Error - 404)
```json
{
  "success": false,
  "message": "Kullanıcı bulunamadı",
  "data": null,
  "error": {
    "code": "NOT_FOUND",
    "message": "Kullanıcı bulunamadı",
    "details": "sql: no rows in result set"
  },
  "timestamp": "string"
}
```

---

## DELETE /api/users/{id}
UUID'ye göre bir kullanıcıyı siler.

### Path Parameters
- **id**: Kullanıcı UUID'si

### Response Body (Success - 200)
```json
{
  "success": true,
  "message": "Kullanıcı başarıyla silindi",
  "data": null,
  "error": null,
  "timestamp": "string"
}
```
