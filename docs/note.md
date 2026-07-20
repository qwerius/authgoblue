```text
Request
   |
   v
Extract Token
   |
   +--> Authorization Header
   |
   +--> Cookie (jika aktif)
   |
   v
Parse & Validate JWT
   |
   v
Check Session & Revoke
   |
   v
Set Claims Context
   |
   v
Controller
```