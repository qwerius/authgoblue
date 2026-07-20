# Panduan penggunaan AuthGoBlue

Dokumen ini menjelaskan cara mengintegrasikan module `authgoblue` ke aplikasi Go yang menggunakan Fiber.

## 1. Instalasi

```bash
go get authgoblue
```

## 2. Inisialisasi module

```go
package main

import (
    "time"

    "authgoblue"
)

func main() {
    agb := authgoblue.New(authgoblue.Config{
        Secret:          "super-secret-key",
        Issuer:          "example-service",
        AccessTokenTTL:  15 * time.Minute,
        RefreshTokenTTL: 7 * 24 * time.Hour,
        Header:          "Authorization",
        Prefix:          "Bearer",
        Cookie:          false,
        CookieName:      "authgoblue_token",
    })

    _ = agb
}
```

## 3. Generate token

```go
accessToken, err := agb.Token.GenerateAccessToken(authgoblue.Claims{
    UserID:      "user-123",
    Username:    "alice",
    Email:       "alice@example.com",
    Role:        "admin",
    Permissions: []string{"read", "write"},
})
```

Catatan:
- `GenerateAccessToken` menghasilkan JWT untuk request yang memerlukan autentikasi.
- `GenerateRefreshToken` digunakan untuk refresh session.

## 4. Middleware autentikasi

```go
app.Use("/protected", agb.Middleware.RequireAuth())
```

Dengan middleware ini, request menuju `/protected` akan memeriksa token pada header `Authorization`.

## 5. Route yang butuh role atau permission

```go
app.Use("/admin", agb.Middleware.RequireAuth())
app.Use("/admin", agb.Middleware.RequireRole("admin"))

app.Use("/reports", agb.Middleware.RequireAuth())
app.Use("/reports", agb.Middleware.RequirePermission("read"))
```

## 6. Mengambil claims dari konteks request

```go
userID, err := agb.Context.UserID(c)
role, err := agb.Context.Role(c)
permissions, err := agb.Context.Permissions(c)
```

## 7. Contoh lengkap

Lihat folder `examples/basic` untuk contoh minimal dan `examples/role_permission` untuk contoh route protected dengan role.
