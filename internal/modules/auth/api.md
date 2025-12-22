# Auth Module API Documentation

## POST /login
Kullanıcı girişi yapar ve JWT token döner.

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

### JWT Token Payload
Token decode edildiğinde aşağıdaki claims içerir:
```json
{
  "user_id": "uuid",
  "username": "string",
  "role": "string",
  "exp": 1234567890
}
```

### Validation Rules
- **username**: Zorunlu (required)
- **password**: Zorunlu (required)

### Test Kullanıcıları
Geliştirme ortamında aşağıdaki test kullanıcıları kullanılabilir:
- `admin` / `123` - ADMIN rolü
- `sekreter` / `123` - SEKRETER rolü
