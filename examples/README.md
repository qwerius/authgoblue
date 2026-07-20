# Contoh penggunaan

Folder `examples/basic` berisi contoh aplikasi Fiber minimal yang menggunakan `authgoblue`.

## Menjalankan contoh

```bash
go run ./examples/basic
```

## Alur contoh

1. `GET /login` akan menghasilkan `access_token` dan `refresh_token`.
2. `GET /protected/me` memerlukan token lewat header `Authorization: Bearer <token>`.
3. Middleware `RequireAuth()` akan memparsing dan memvalidasi access token.

## Contoh request

```bash
curl http://localhost:3000/login
curl http://localhost:3000/protected/me \
  -H "Authorization: Bearer <access_token>"
```
