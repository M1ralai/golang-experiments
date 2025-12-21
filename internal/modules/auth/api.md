# Auth Module API Documentation

## POST /login
Kullanıcı girişi yapar.

### Request Body
```json
{
  "username": "string", // Zorunlu
  "password": "string"  // Zorunlu
}
```

### Response Body (Success - 200)
```json
{
  "success": true,
  "message": "Giriş başarılı",
  "data": {
    "token": "string",
    "username": "string",
    "role": "string",
    "ad": "string",
    "soyad": "string"
  },
  "error": null,
  "timestamp": "string"
}
```

### Validation Rules
- **username**: Zorunlu (required)
- **password**: Zorunlu (required)
