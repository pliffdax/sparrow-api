# Sparrow API

## Functionality

### Public endpoints
- `GET /health` — returns JSON with status, current time and API version
- `POST /auth/register` — registers a new user  
  Request: `{ "name": "Alice", "password": "123" }` → 201
- `POST /auth/login` — authenticates and returns JWT token  
  Response: `{ "access_token": "<jwt>", "token_type": "Bearer" }`

### Protected endpoints  
Require header: `Authorization: Bearer <jwt>`

#### Users
- `GET /users` — list users → 200
- `GET /users/{id}` — get user by ID → 200/404
- `DELETE /users/{id}` — delete user by ID → 204/404

#### Categories
- `POST /categories` — create a category `{ "title": "Food" }` → 201
- `GET /categories` — list categories → 200
- `DELETE /categories/{id}` — delete category → 204/404

#### Records
- `POST /records` — create a record  
  `{ "user_id": 1, "category_id": 1, "amount": 99.5, "created_at"?: RFC3339 }`  
  Validates user + category existence → 201/400
- `GET /records?user_id=&category_id=` — filter by one or both params (at least one required) → 200/400
- `GET /records/{id}` — get record by ID → 200/404
- `DELETE /records/{id}` — delete record → 204/404


## Project structure
```text
sparrow-api
.
├── cmd/
│   └── server/
│       └── main.go
├── docker-compose.yml
├── Dockerfile
├── docs/
│   ├── lab1/
│   ├── lab2/
│   ├── lab3/
│   └── lab4/
├── go.mod
├── go.sum
├── internal/
│   ├── app/
│   │   └── app.go
│   ├── domain/
│   │   ├── category.go
│   │   ├── record.go
│   │   └── user.go
│   ├── http/
│   │   ├── handlers/
│   │   │   ├── auth.go
│   │   │   ├── category.go
│   │   │   ├── health.go
│   │   │   ├── record.go
│   │   │   └── user.go
│   │   ├── middleware/
│   │   │   └── auth.go
│   │   └── router.go
│   ├── security/
│   │   └── jwt.go
│   ├── storage/
│   │   ├── memory/
│   │   │   ├── categories.go
│   │   │   ├── records.go
│   │   │   └── users.go
│   │   ├── postgres/
│   │   │   ├── categories.go
│   │   │   ├── db.go
│   │   │   ├── records.go
│   │   │   └── users.go
│   │   └── storage.go
│   └── util/
│       ├── env.go
│       └── json.go
├── migrations/
│   ├── 0001_init.sql
│   ├── 0002_seed.sql
│   └── 0003_add_password_hash.sql
└── README.md
```

### Environment variables required

| Variable      | Description |
|--------------|-------------|
| POSTGRES_DB | database name |
| POSTGRES_USER | database user |
| POSTGRES_PASSWORD | database password |
| DB_DSN | connection string for API |
| JWT_SECRET | secret key used to sign JWT tokens |
| JWT_TTL | token lifetime (e.g. `24h`) |

Example:
```env
POSTGRES_DB=sparrow
POSTGRES_USER=sparrow
POSTGRES_PASSWORD=secret
DB_DSN=postgres://sparrow:secret@db:5432/sparrow?sslmode=disable
JWT_SECRET=some_super_secret_key
JWT_TTL=24h
PORT=8080
```

## How to start localy
```bash
go run ./cmd/server
```

## Run with Docker
```bash
# build image
docker build -t sparrow-api:latest .

# run container
docker run --rm -p 8080:8080 -e PORT=8080 sparrow-api:latest
```

## Run with Docker Compose
```bash
docker-compose up --build
```

## Database Migrations

PostgreSQL migrations are stored in the `migrations/` folder and must be applied
before using authenticated endpoints. There are two options:

---

### Option A — Fresh start (recommended for local development)

If you want PostgreSQL to apply migrations automatically on startup:

```bash
docker compose down -v
docker compose up --build
```

The SQL files inside ./migrations will run automatically only on the first launch
of a new Postgres volume.

---

### Option B — Apply migrations manually (if DB already exists)

Check existing tables:

```bash
docker exec -it sparrow-db \
  psql -U $POSTGRES_USER -d $POSTGRES_DB -c '\dt'
```

Apply migrations one-by-one:

```bash
docker exec -it sparrow-db \
  psql -U $POSTGRES_USER -d $POSTGRES_DB \
  -f /docker-entrypoint-initdb.d/0001_init.sql

docker exec -it sparrow-db \
  psql -U $POSTGRES_USER -d $POSTGRES_DB \
  -f /docker-entrypoint-initdb.d/0002_seed.sql

docker exec -it sparrow-db \
  psql -U $POSTGRES_USER -d $POSTGRES_DB \
  -f /docker-entrypoint-initdb.d/0003_add_password_hash.sql
```

> Use this only if your database volume already exists and tables were not created automatically.

## Deploy
The application is deployed on Render.

Accessible at: https://sparrow-api-l8pp.onrender.com

## Git

Using Conventional Commits style.

Examples:
```text
feat: add chi router with /health endpoint
chore(docker): add Dockerfile and docker-compose setup
```

## Variant

**Group:** IO-35  
**Variant formula:** 35 % 3 = 2  
**Resulting task:** **Користувацькі категорії витрат**

Implemented according to the variant:
- Added global (`is_global = true`) and user (`is_global = false, owner_id = user_id`) expense categories.
- Global categories are visible to all users.
- User categories are visible only to their creators.
- Access control implemented: users cannot delete or use someone else's categories (HTTP 403).
- All logic verified via Postman tests in `docs/lab3/`.

## Author
- GitHub: [@Pliffdax](https://github.com/Pliffdax)  
- Telegram: [@Pliffdax](https://t.me/Pliffdax)

---

© KPI — Pliffdax – Sparrow API — 2025
