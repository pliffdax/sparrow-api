# Sparrow API

## Functionality
- `GET /health` — returns JSON with ok status, current time and version of the project
- `POST /users` — create a user { "name": "Alice" } → 201
- `GET /users` — list users → 200
- `GET /users/{id}` — get user by id → 200/404
- `DELETE /users/{id}` — delete user → 204/404
- `POST /categories` — create a category { "title": "Food" } → 201
- `GET /categories` — list categories → 200
- `DELETE /categories/{id}` — delete category → 204/404
- `POST /records` — create a record { "user_id": 1, "category_id": 1, "amount": 99.5, "created_at"?: RFC3339 } Validates that user_id and category_id exist. → 201/400
- `GET /records?user_id=&category_id=` — filter by one or both params (at least one required) → 200/400
- `GET /records/{id}` — get record by id → 200/404
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
│   │   │   ├── category.go
│   │   │   ├── health.go
│   │   │   ├── record.go
│   │   │   └── user.go
│   │   └── router.go
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
│   └── 0002_seed.sql
└── README.md
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
