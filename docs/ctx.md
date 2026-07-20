# AuthGoBlue Context (`ctx`)

Package `ctx` menyediakan helper untuk mengakses informasi autentikasi yang sudah tersedia pada request context.

`ctx.Service` digunakan setelah middleware AuthGoBlue berhasil melakukan:

1. Membaca token dari request
2. Memvalidasi token
3. Membuat `claims.Claims`
4. Menyimpan informasi autentikasi ke `fiber.Ctx`