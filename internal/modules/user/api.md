# User Module API Documentation

## GET /api/users
Tüm kayıtlı kullanıcıları listeler.

### Response Body (Success - 200)
```json
[
  {
    "id": 1,
    "username": "string",
    "role": "string",
    "ad": "string",
    "soyad": "string",
    "telefon": "string",
    "email": "string"
  }
]
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
    "id": 1,
    "username": "string",
    "role": "string",
    "ad": "string",
    "soyad": "string",
    "telefon": "string",
    "email": "string"
  },
  "error": null
}
```

### Validation Rules
- **username**: Zorunlu (required)
- **password**: Zorunlu (required) ve en az 3 karakter
- **role**: Zorunlu (required)

---

## DELETE /api/users/{id}
ID'ye göre bir kullanıcıyı siler.

### Response Body (Success - 200)
```json
{
  "success": true,
  "message": "Kullanıcı başarıyla silindi",
  "data": null
}
```
