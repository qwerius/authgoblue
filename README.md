````md
# AuthGoBlue

Authentication and Session Management module for Go Fiber applications using JWT.

AuthGoBlue adalah library authentication reusable untuk aplikasi Go Fiber yang menyediakan sistem autentikasi lengkap berbasis JWT dengan dukungan session management, refresh token security, revoke system, dan authorization.


## Features

AuthGoBlue menyediakan:

- JWT Access Token
- JWT Refresh Token
- JWT signature validation
- Token expiration validation
- Token type validation
- Refresh token rotation
- Refresh token reuse detection
- Authentication middleware
- Claims management
- Context helper
- Header Bearer authentication
- Cookie authentication support
- Session management
- Memory session store
- Redis session store
- Custom session store
- Token revoke
- Session revoke
- Role authorization
- Permission authorization
- Multiple device session
- Maximum session limit


---

# Installation

Install AuthGoBlue:

```bash
go get github.com/blueink/authgoblue
````

Import package:

```go
import (
	"github.com/blueink/authgoblue"
	"github.com/blueink/authgoblue/claims"
)
```

---

# Quick Start

Initialize AuthGoBlue:

```go
agb := authgoblue.New(
	authgoblue.Config{

		Secret:
			"super-secret-key",

		Issuer:
			"example-service",

		AccessTokenTTL:
			15 * time.Minute,

		RefreshTokenTTL:
			7 * 24 * time.Hour,

		Header:
			"Authorization",

		Prefix:
			"Bearer",

		Cookie:
			false,

	},
)
```

---

# Configuration

## Secret

Secret digunakan untuk signing dan verification JWT.

Example:

```go
Secret: "super-secret-key"
```

Production:

Gunakan environment variable.

```bash
AUTH_SECRET=your-secret-key
```

---

## Issuer

Issuer adalah identitas service pembuat token.

Example:

```go
Issuer:"auth-service"
```

---

## AccessTokenTTL

Masa berlaku access token.

Example:

```go
AccessTokenTTL:
	15 * time.Minute
```

Rekomendasi:

```
5 - 30 menit
```

---

## RefreshTokenTTL

Masa berlaku refresh token.

Example:

```go
RefreshTokenTTL:
	7 * 24 * time.Hour
```

---

## Header

Header tempat token dikirim.

Default:

```go
Authorization
```

Request:

```
Authorization: Bearer ACCESS_TOKEN
```

---

## Prefix

Prefix token.

Default:

```go
Bearer
```

Example:

```
Bearer eyJhbGciOiJIUzI1...
```

---

## Cookie

Mode penyimpanan token.

Header:

```go
Cookie:false
```

Cookie:

```go
Cookie:true
```

---

# Authentication Flow

```
Client

 |

Login

 |

Generate Token

 |

Receive Access Token

 |

API Request

 |

Authorization Header

 |

AuthGoBlue Middleware

 |

Parse JWT

 |

Validate JWT

 |

Session Check

 |

Set Claims Context

 |

Handler
```

---

# JWT Token

AuthGoBlue menggunakan dua jenis token.

## Access Token

Access token digunakan untuk request API.

Karakteristik:

* Lifetime pendek
* Digunakan pada API request
* Signature validation
* Expiration validation

Example:

```go
token, err :=
agb.Token.GenerateAccessToken(
	claims.Claims{

		UserID:
			"user-123",

		Role:
			"admin",

	},
)
```

---

## Refresh Token

Refresh token digunakan untuk mendapatkan access token baru.

Support:

* Rotation
* Reuse detection
* Session binding

Example:

```go
refresh, err :=
agb.Token.GenerateRefreshToken(
	claims.Claims{

		UserID:
			"user-123",

	},
)
```

---

# Claims

Claims adalah informasi user yang disimpan dalam JWT.

Example:

```go
claims.Claims{

	UserID:
		"user-123",

	Username:
		"alice",

	Email:
		"alice@example.com",

	Role:
		"admin",

	Permissions:
		[]string{
			"read",
			"write",
		},

}
```

Field:

| Field       | Description               |
| ----------- | ------------------------- |
| UserID      | User identifier           |
| Username    | Username                  |
| Email       | Email user                |
| Role        | User role                 |
| Permissions | Permission list           |
| SessionID   | Active session identifier |
| TokenType   | Access / Refresh          |

---

# Authentication Middleware

Protect route:

```go
app.Use(
	"/protected",
	agb.Middleware.RequireAuth(),
)
```

Example:

```go
app.Get(
	"/protected/me",
	func(c fiber.Ctx) error {

		return c.JSON(
			fiber.Map{
				"message":
					"authenticated",
			},
		)

	},
)
```

---

# Middleware Options

AuthGoBlue menyediakan custom middleware behavior.

Example:

```go
agb.Middleware.RequireAuthWith(
	middleware.Options{

		SkipSessionCheck:
			true,

	},
)
```

Available options:

## SkipValidation

```go
SkipValidation:true
```

Hanya melakukan parsing token.

Tidak melakukan:

* expiration validation
* token validation

---

## SkipSessionCheck

```go
SkipSessionCheck:true
```

Mode JWT Stateless.

Behavior:

* Tidak melakukan SessionStore lookup
* Tidak menggunakan Redis
* Tidak melakukan session verification

Cocok untuk:

* Internal service
* High throughput API
* Stateless authentication

---

## SkipRevokeCheck

```go
SkipRevokeCheck:true
```

Tidak melakukan pengecekan revoke.

---

# Session Management

AuthGoBlue menggunakan SessionID untuk tracking session.

Session identifier:

```
SessionID
```

Tidak menggunakan JTI sebagai session identifier.

Session structure:

```
SessionID

UserID

CreatedAt

ExpiresAt

Revoked
```

---

# Session Store

AuthGoBlue support beberapa storage.

## Memory Session Store

Default store.

Cocok untuk:

* Development
* Testing
* Single instance application

Performance:

```
Create:
~600 ns


Get:
~37 ns
```

---

## Redis Session Store

Redis digunakan untuk production distributed system.

Cocok untuk:

* Multiple server
* Load balancing
* Shared session state

Performance:

```
Create:
~520 µs


Get:
~270 µs
```

---

# Multiple Device Session

AuthGoBlue mendukung login banyak device.

Example:

```
User

 |

+----------------+

Device A

Session ID: xxx


+----------------+

Device B

Session ID: yyy

```

Setiap device mempunyai session berbeda.

---

# Session Limit

Support:

* Maximum active session
* Automatic oldest session revoke

Example:

```
MaxSession = 3


Login ke-4

|

Session lama dihapus
```

---

# Logout

## Logout Current Device

Hanya revoke session aktif.

```
Device A

Logout

Session A revoked
```

---

## Logout All Devices

Semua session user dapat direvoke.

```
Device A
Device B
Device C

Logout All

|

All session revoked
```

---

# Revoke System

AuthGoBlue mempunyai revoke store.

Support:

* Memory revoke store
* Redis revoke store
* Custom revoke store

Flow:

```
Request

 |

JWT

 |

SessionID

 |

Revoke Store

 |

Allow / Reject
```

---

# Refresh Token Security

## Rotation

Flow:

```
Refresh Token A

       |

       v

Refresh Token B
```

Token lama menjadi invalid.

---

## Reuse Detection

Jika refresh token lama digunakan kembali:

```
Old Refresh Token

        |

        v

Rejected
```

---

# Authorization

## Role Authorization

Example:

```go
agb.Middleware.RequireRole(
	"admin",
)
```

---

## Permission Authorization

Example:

```go
agb.Middleware.RequirePermission(
	"user.delete",
)
```

---

# Context Helper

Mengambil informasi user dari Fiber context.

User ID:

```go
userID, err :=
agb.Context.UserID(c)
```

Role:

```go
role, err :=
agb.Context.Role(c)
```

Permission:

```go
permissions, err :=
agb.Context.Permissions(c)
```

---

# Benchmark

Hardware:

```
CPU:
Intel Core i9-9900K

OS:
Windows AMD64
```

## JWT Performance

| Operation             | Result |
| --------------------- | -----: |
| Generate Access Token |  ~3 µs |
| Parse Access Token    |  ~5 µs |
| Validate Access Token | ~20 ns |

---

## Middleware Performance

| Mode          |  Result |
| ------------- | ------: |
| JWT + Session |  ~28 µs |
| JWT Stateless |  ~27 µs |
| Redis Session | ~310 µs |

---

## Session Performance

| Store  |  Create |     Get |
| ------ | ------: | ------: |
| Memory | ~600 ns |  ~37 ns |
| Redis  | ~520 µs | ~270 µs |

---

# Production Recommendation

Gunakan:

* HTTPS
* Strong secret
* Environment variable
* Short access token TTL
* Refresh token rotation
* HttpOnly cookie untuk browser
* Redis session untuk multi server

---

# Roadmap

Future development:

* OAuth Provider
* JWT key rotation
* Database session adapter
* Audit logging
* Device fingerprinting
* Rate limit integration

---

# License

MIT License

```
```
