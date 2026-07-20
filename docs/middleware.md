# AuthGoBlue Middleware

Middleware AuthGoBlue digunakan untuk melindungi route yang membutuhkan autentikasi.

Middleware bertugas untuk:

1. Mengambil token dari request
2. Membaca JWT access token
3. Memvalidasi token
4. Memeriksa session aktif
5. Memeriksa revoke/blacklist
6. Menyimpan claims authentication ke context request

Setelah middleware berhasil, informasi user dapat diakses melalui facade AuthGoBlue:

```go
agb.UserID(c)
agb.Email(c)
agb.Role(c)
agb.Permissions(c)
```
<img src="images/authrequired.png" alt="Auth Required Flow" width="600">
---

## Initialization

Middleware sudah tersedia otomatis setelah membuat instance AuthGoBlue.

Contoh:

```go
agb := authgoblue.New(
	authgoblue.Config{
		Secret: "secret-key",
		Issuer: "my-app",
	},
)

// Access middleware
agb.Middleware
```

---

## Configuration



```go
authgoblue.Config
```

Developer tidak perlu membuat konfigurasi middleware secara terpisah.

Contoh:

```go
agb := authgoblue.New(
	authgoblue.Config{
		Secret: "secret-key",
		Issuer: "my-app",

		Header: "Authorization",
		Prefix: "Bearer",

		Cookie: true,
		CookieName: "auth_token",
	},
)
```

---

### Middleware Configuration

#### Header

Menentukan nama HTTP header yang digunakan untuk mengambil token.

Default:

```go
Header: "Authorization"
```

Contoh request:

```http
Authorization: Bearer eyJhbGciOi...
```

---

#### Prefix

Menentukan prefix token pada authorization header.

Default:

```go
Prefix: "Bearer"
```

Contoh:

```http
Authorization: Bearer eyJhbGciOi...
```

Middleware akan mengambil nilai setelah prefix:

```
eyJhbGciOi...
```

---

#### Cookie

Mengaktifkan pengambilan token melalui cookie.

Default:

```go
Cookie: false
```

Contoh:

```go
Cookie: true
```

Jika aktif, middleware dapat membaca token dari cookie.

---

#### CookieName

Menentukan nama cookie yang berisi token.

Default:

```go
CookieName: "github.com/qwerius/authgoblue_token"
```

Contoh:

```go
CookieName: "auth_token"
```

Request:

```http
Cookie: auth_token=eyJhbGciOi...
```

---

### Default Configuration

AuthGoBlue menyediakan konfigurasi default:

```go
authgoblue.DefaultConfig()
```

Hasil default:

```go
authgoblue.Config{
	AccessTokenTTL: 15 * time.Minute,

	RefreshTokenTTL: 7 * 24 * time.Hour,

	Header: "Authorization",

	Prefix: "Bearer",

	Cookie: false,

	CookieName: "github.com/qwerius/authgoblue_token",

	MaxSessions: 5,
}
```


### RequireAuth()

`RequireAuth()` adalah middleware default untuk route yang membutuhkan login.

Signature:

```go
func (s *Service) RequireAuth() fiber.Handler
```

Contoh:

```go
app.Get(
	"/profile",
	agb.Middleware.RequireAuth(),
	Profile,
)


Flow:

```
```
Request
   |
   v
Extract Token
   |
   v
Parse Access Token
   |
   v
Validate Token
   |
   v
Check Session
   |
   v
Check Revoke
   |
   v
Set Authentication Context
   |
   v
Controller
```

---

# Token Source

Sumber token ditentukan melalui konfigurasi AuthGoBlue.

Contoh:

```go
authgoblue.Config{
	Header: "Authorization",
	Prefix: "Bearer",

	Cookie: true,
	CookieName: "auth_token",
}
```

Urutan pembacaan token:

```
Authorization Header
        |
        v
Cookie (jika Cookie=true)
```

---

# Authorization Header

Default:

```go
authgoblue.Config{
	Header: "Authorization",
	Prefix: "Bearer",
}
```

Request:

```http
Authorization: Bearer eyJhbGciOi...
```

Middleware mengambil:

```
eyJhbGciOi...
```

---

# Cookie Authentication

Aktifkan cookie:

```go
authgoblue.Config{
	Cookie: true,
	CookieName: "auth_token",
}
```

Request:

```http
Cookie: auth_token=eyJhbGciOi...
```

Middleware akan membaca token dari cookie apabila header tidak tersedia.

---

# RequireAuthWith()

Untuk kebutuhan khusus, gunakan:

```go
func (s *Service) RequireAuthWith(
	options Options,
) fiber.Handler
```

Contoh:

```go
app.Get(
	"/internal",
	agb.Middleware.RequireAuthWith(
		middleware.Options{
			SkipSessionCheck: true,
		},
	),
	Handler,
)
```

---

# Options

## SkipValidation

Default:

```go
false
```

Jika `true`, middleware melewati validasi access token.

Contoh:

```go
middleware.Options{
	SkipValidation: true,
}
```

Perilaku:

```
Parse Token
    |
    v
Skip Validation
```

Gunakan hanya untuk kebutuhan khusus.

---

## SkipSessionCheck

Default:

```go
false
```

Jika `true`, middleware tidak melakukan pengecekan session.

Contoh:

```go
middleware.Options{
	SkipSessionCheck: true,
}
```

Cocok untuk mode JWT stateless:

```
JWT valid
Signature valid
Expiration valid

Tidak ada SessionStore lookup
```

---

## SkipRevokeCheck

Default:

```go
false
```

Jika `true`, middleware melewati pengecekan revoke.

Contoh:

```go
middleware.Options{
	SkipRevokeCheck: true,
}
```

---

# Controller Usage

Setelah middleware berhasil, developer dapat mengambil informasi user melalui facade AuthGoBlue.

Contoh:

```go
func Profile(
	c fiber.Ctx,
) error {

	userID, err := agb.UserID(c)

	if err != nil {
		return err
	}

	email, err := agb.Email(c)

	if err != nil {
		return err
	}

	return c.JSON(
		fiber.Map{
			"id": userID,
			"email": email,
		},
	)
}
```

---

# Protected Route Example

Contoh lengkap:

```go
package main

import (
	"github.com/qwerius/authgoblue"

	"github.com/gofiber/fiber/v3"
)

func main() {

	app := fiber.New()

	agb := authgoblue.New(
		authgoblue.Config{
			Secret: "secret",
			Issuer: "application",
		},
	)

	app.Get(
		"/me",
		agb.Middleware.RequireAuth(),
		func(c fiber.Ctx) error {

			email, err := agb.Email(c)

			if err != nil {
				return err
			}

			return c.JSON(
				fiber.Map{
					"email": email,
				},
			)
		},
	)

	app.Listen(":3000")
}
```

---

# Security Behavior

Default middleware:

```
✓ Token extraction
✓ JWT parsing
✓ Signature validation
✓ Issuer validation
✓ Expiration validation
✓ Session validation
✓ Revoke validation
✓ Claims storage
```

---

# Recommendation

Gunakan:

```go
agb.Middleware.RequireAuth()
```

untuk hampir semua route.

Gunakan:

```go
agb.Middleware.RequireAuthWith(...)
```

hanya jika membutuhkan perilaku khusus seperti:

- JWT stateless
- internal service
- kebutuhan testing
- bypass tertentu

