# Personal website

- Stack: Htmx and GO
- Blog: GO custom markdown to HTML conversion

## Env vars

The app (`cmd/web`) and migrations job (`cmd/migrate`) need these set
(local dev loads them from `.local.env` via `docker compose up`):

| Var | Required | Format / example |
| --- | --- | --- |
| `DB_HOST` | yes | e.g. `db` (compose) / `localhost` (local tests) |
| `DB_PORT` | yes | e.g. `5432` |
| `DB_USER` | yes | e.g. `postgres` |
| `DB_PASSWORD` | yes | e.g. `docker` |
| `DB_NAME` | yes | e.g. `postgres` |
| `RSA_KEY` | yes for login/session | base64 (std encoding, single line) of a PEM `RSA PRIVATE KEY` (PKCS#1). Used to sign/verify RS256 session JWTs. |
| `ADMIN_HTPASSWD` | yes for login | `user:bcrypt-hash`, e.g. `sam:$2y$05$...`. Compared with bcrypt in `POST /auth`. |

Generate the auth vars:

```bash
# RSA private key -> single-line base64 for RSA_KEY
openssl genrsa -traditional -out private.pem 2048
base64 -w0 private.pem  # macOS: base64 -b0 private.pem

# bcrypt htpasswd entry for ADMIN_HTPASSWD (needs apache2-utils)
htpasswd -nbB sam <password>
```
