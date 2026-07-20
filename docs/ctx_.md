# AuthGoBlue Context (`ctx`)

Package `ctx` menyediakan helper untuk mengakses informasi autentikasi yang sudah tersedia pada request context.

`ctx.Service` digunakan setelah middleware AuthGoBlue berhasil melakukan:

1. Membaca token dari request
2. Memvalidasi token
3. Membuat `claims.Claims`
4. Menyimpan informasi autentikasi ke `fiber.Ctx`

Developer tidak perlu melakukan parsing JWT secara manual.

---

# Initialization

`ctx.Service` dibuat otomatis saat membuat instance AuthGoBlue.

```go
auth := authgoblue.New(config)
```

Service dapat diakses melalui:

```go
auth.Context
```

Contoh:

```go
userID, err := auth.Context.UserID(c)
```

---

# Service Structure

```go
type Service struct{}
```

Constructor:

```go
func NewService() *Service
```

---

# Claims

Semua helper identity pada package `ctx` mengambil data dari:

```go
authClaims, err := s.Claims(c)
```

Claims berisi:

```go
type Claims struct {
	UserID    string
	SessionID string
	Username  string
	Email     string

	Role string

	Permissions []string

	TokenType TokenType

	Issuer string

	IssuedAt  int64
	ExpiresAt int64
}
```

---

# User Information

## UserID()

Mengambil ID user yang sedang login.

Signature:

```go
func (s *Service) UserID(
	c fiber.Ctx,
) (string, error)
```

Contoh:

```go
id, err := auth.Context.UserID(c)

if err != nil {
	return err
}
```

Return:

```
user_id
```

---

## Username()

Mengambil username user.

Signature:

```go
func (s *Service) Username(
	c fiber.Ctx,
) (string, error)
```

Contoh:

```go
username, err := auth.Context.Username(c)
```

Return:

```
username
```

---

## Email()

Mengambil email user.

Signature:

```go
func (s *Service) Email(
	c fiber.Ctx,
) (string, error)
```

Contoh:

```go
email, err := auth.Context.Email(c)
```

Return:

```
email@example.com
```

---

# Role

## Role()

Mengambil role user.

Signature:

```go
func (s *Service) Role(
	c fiber.Ctx,
) (string, error)
```

Contoh:

```go
role, err := auth.Context.Role(c)
```

Contoh nilai:

```
admin
user
moderator
```

---

# Permission

## Permissions()

Mengambil seluruh permission user.

Signature:

```go
func (s *Service) Permissions(
	c fiber.Ctx,
) ([]string, error)
```

Contoh:

```go
permissions, err := auth.Context.Permissions(c)
```

Return:

```go
[]string{
	"users:create",
	"users:update",
}
```

---

## HasPermission()

Mengecek permission tertentu.

Signature:

```go
func (s *Service) HasPermission(
	c fiber.Ctx,
	permission string,
) (bool, error)
```

Contoh:

```go
allowed, err := auth.Context.HasPermission(
	c,
	"users:create",
)

if !allowed {
	return fiber.ErrForbidden
}
```

Return:

```
true
false
```

---

# Session Context

Package `ctx` juga menyediakan akses terhadap session aktif.

Session disimpan menggunakan:

```go
c.Locals()
```

dengan key:

```
authgoblue_session
```

---

## SetSession()

Menyimpan session ke request context.

Signature:

```go
func (s *Service) SetSession(
	c fiber.Ctx,
	value session.Session,
)
```

Contoh:

```go
auth.Context.SetSession(
	c,
	currentSession,
)
```

---

## Session()

Mengambil session aktif.

Signature:

```go
func (s *Service) Session(
	c fiber.Ctx,
) (session.Session, error)
```

Contoh:

```go
sess, err := auth.Context.Session(c)
```

Return:

```go
session.Session
```

---

## SessionID()

Mengambil ID session aktif.

Signature:

```go
func (s *Service) SessionID(
	c fiber.Ctx,
) (string, error)
```

Contoh:

```go
sessionID, err := auth.Context.SessionID(c)
```

Return:

```
session_id
```

---

# Token Type

Package `ctx` menyediakan helper untuk membaca tipe token dari Claims.

Token type berasal dari:

```go
claims.TokenType
```

Dengan nilai:

```go
claims.TokenTypeAccess

claims.TokenTypeRefresh
```

---

## TokenType()

Mengambil tipe token.

Signature:

```go
func (s *Service) TokenType(
	c fiber.Ctx,
) (claims.TokenType, error)
```

Contoh:

```go
tokenType, err := auth.Context.TokenType(c)
```

---

## IsAccessToken()

Mengecek apakah token adalah access token.

Signature:

```go
func (s *Service) IsAccessToken(
	c fiber.Ctx,
) (bool, error)
```

Contoh:

```go
ok, err := auth.Context.IsAccessToken(c)

if ok {
	// access token
}
```

---

## IsRefreshToken()

Mengecek apakah token adalah refresh token.

Signature:

```go
func (s *Service) IsRefreshToken(
	c fiber.Ctx,
) (bool, error)
```

Contoh:

```go
ok, err := auth.Context.IsRefreshToken(c)

if ok {
	// refresh token
}
```

---

# Error Handling

Semua method yang membaca data authentication mengembalikan error.

Contoh:

```go
email, err := auth.Context.Email(c)

if err != nil {
	return err
}
```

Error dapat terjadi ketika:

* Claims tidak ditemukan
* Request belum melewati middleware authentication
* Session belum tersedia

---

# Controller Example

Contoh penggunaan lengkap:

```go
func Profile(c fiber.Ctx) error {

	userID, err := auth.Context.UserID(c)

	if err != nil {
		return err
	}

	email, err := auth.Context.Email(c)

	if err != nil {
		return err
	}

	role, err := auth.Context.Role(c)

	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"id":    userID,
		"email": email,
		"role":  role,
	})
}
```

---

# API Summary

| Method             | Description                    |
| ------------------ | ------------------------------ |
| `UserID()`         | Mendapatkan user ID            |
| `Username()`       | Mendapatkan username           |
| `Email()`          | Mendapatkan email              |
| `Role()`           | Mendapatkan role               |
| `Permissions()`    | Mendapatkan daftar permission  |
| `HasPermission()`  | Mengecek permission            |
| `SetSession()`     | Menyimpan session pada context |
| `Session()`        | Mengambil session aktif        |
| `SessionID()`      | Mengambil ID session           |
| `TokenType()`      | Mengambil tipe token           |
| `IsAccessToken()`  | Mengecek access token          |
| `IsRefreshToken()` | Mengecek refresh token         |

---

# Responsibility

Package `ctx` bertanggung jawab untuk:

* Membaca authentication context
* Mengambil Claims
* Mengambil user identity
* Mengambil authorization data
* Mengakses session aktif

Package `ctx` tidak bertanggung jawab untuk:

* Membuat JWT
* Validasi JWT
* Login
* Logout
* Refresh token
* Hash password
* Penyimpanan session
