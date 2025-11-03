
# Лабораторна робота №4 — JWT‑авторизація

> Продовження ЛР‑3. У цій роботі додаємо реєстрацію/вхід користувача, видачу JWT‑токену та захист ресурсів API.

---

## 1. Мета роботи
- Додати до REST‑API підтримку аутентифікації на основі **JWT**.
- Розмежувати публічні та захищені ендпоінти.
- Зберігати хеш пароля в БД, не передаючи/не зберігаючи пароль у відкритому вигляді.
- Перевірити роботу локально (Docker + Postman/cURL) і підготувати артефакти для звіту.

---

## 2. Короткий опис змін
- **Нові ендпоінти (публічні):**
  - `POST /auth/register` — створення користувача з паролем.
  - `POST /auth/login` — вхід, повертає `access_token` (**JWT**, `Bearer`).

- **JWT‑middleware**: захищає всі інші маршрути (`/users`, `/categories`, `/records`).

- **Міграція БД**: додано поле `password_hash` у таблицю `users`:
  - `migrations/0003_add_password_hash.sql`

- **Конфіг (env)**: додано змінні
  - `JWT_SECRET` — секрет підпису.
  - `JWT_TTL` — час життя токена (наприклад, `24h`).

- **Storage**:
  - новий інтерфейс `AuthUserStore` (пакет `internal/storage`):  
    ```go
    type AuthUserStore interface {
        CreateWithPassword(name, passwordHash string) (domain.User, error)
        FindAuth(name string) (id int64, passwordHash string, err error)
    }
    ```
  - реалізація у `internal/storage/postgres/users.go`:
    - `CreateWithPassword(name, passwordHash string) (domain.User, error)`
    - `FindAuth(name string) (int64, string, error)`

- **Router**: публічні `/auth/*` + група закритих шляхів під `AuthRequired` middleware.

---

## 3. Залежності та змінні середовища

### 3.1. Залежності (модуль security)
Файл: `internal/security/jwt.go` — створення/перевірка JWT на основі `JWT_SECRET` і `JWT_TTL` (парсинг тривалості через `time.ParseDuration`).

### 3.2. `.env` (приклад)
```env
POSTGRES_DB=sparrow
POSTGRES_USER=sparrow
POSTGRES_PASSWORD=secret

DB_DSN=postgres://sparrow:secret@db:5432/sparrow?sslmode=disable

JWT_SECRET=some_super_secret_key_!!!_lol
JWT_TTL=24h
PORT=8080
```

---

## 4. Docker‑компоуз (актуальний фрагмент)

```yaml
services:
  db:
    image: postgres:17
    container_name: sparrow-db
    restart: always
    environment:
      POSTGRES_DB: ${POSTGRES_DB}
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data
      # для автозапуску SQL при першій ініціалізації кластера:
      # - ./migrations:/docker-entrypoint-initdb.d
      # якщо кластер вже існує — застосуйте файли вручну (див. нижче)
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER} -d ${POSTGRES_DB}"]
      interval: 3s
      timeout: 3s
      retries: 20

  sparrow-api:
    container_name: sparrow-api
    restart: always
    build:
      context: .
      dockerfile: Dockerfile
    environment:
      PORT: "8080"
      APP_VERSION: "0.1.0"
      DB_DSN: ${DB_DSN}
      JWT_SECRET: ${JWT_SECRET}
      JWT_TTL: ${JWT_TTL}
    depends_on:
      db:
        condition: service_healthy
    ports:
      - "8080:8080"

volumes:
  pgdata:
```

---

## 5. Міграції БД

### Список
- `0001_init.sql` — базові структури (`users`, `categories`, `records`, індекси).
- `0002_seed.sql` — базові записи (опціонально).
- `0003_add_password_hash.sql` — додає поле `password_hash`:

```sql
-- migrations/0003_add_password_hash.sql
ALTER TABLE users ADD COLUMN IF NOT EXISTS password_hash TEXT;
CREATE UNIQUE INDEX IF NOT EXISTS users_name_uq ON users (name);
```

### Застосування (коли кластер вже створено)
```bash
# перевірити таблиці
docker exec -it sparrow-db psql -U sparrow -d sparrow -c '\dt'

# застосувати
docker exec -it sparrow-db psql -U sparrow -d sparrow -f /migrations/0001_init.sql
docker exec -it sparrow-db psql -U sparrow -d sparrow -f /migrations/0002_seed.sql
docker exec -it sparrow-db psql -U sparrow -d sparrow -f /migrations/0003_add_password_hash.sql
```

> Якщо хочете **автозапуск** SQL при першому піднятті контейнера, монтуйте `./migrations` у `/docker-entrypoint-initdb.d` і стовідсотково очистіть volume: `docker compose down -v && docker compose up --build`.

---

## 6. Запуск та перевірка

### 6.1. Підняти середовище
```bash
docker compose down -v
docker compose up --build
```

### 6.2. Перевірити, що API живе
```bash
curl -s http://localhost:8080/health
# очікуємо: OK
```

### 6.3. Тести (cURL)

**Реєстрація**:
```bash
curl -i -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"test","password":"123"}'
```

**Вхід (отримати токен)**:
```bash
curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"name":"test","password":"123"}'
# {"access_token":"<JWT>","token_type":"Bearer"}
```

**Запит до захищеного ресурсу (потрібен Bearer‑токен)**:
```bash
TOKEN="<JWT>"
curl -i -H "Authorization: Bearer $TOKEN" http://localhost:8080/users
```

**Негативні кейси**:
- `GET /users` **без** заголовка `Authorization` → `401 Unauthorized`.
- Невірний пароль у `/auth/login` → `401 Unauthorized`.
- Повторна реєстрація того ж імені → `409 Conflict` (якщо валідація у CreateWithPassword так налаштована).

---

## 7. Контракти ендпоінтів

### 7.1. `POST /auth/register`
**Request**
```json
{
  "name": "alice",
  "password": "secret"
}
```
**Response 201**
```json
{
  "id": 1,
  "name": "alice"
}
```

### 7.2. `POST /auth/login`
**Request**
```json
{
  "name": "alice",
  "password": "secret"
}
```
**Response 200**
```json
{
  "access_token": "<JWT>",
  "token_type": "Bearer"
}
```

### 7.3. Приклади захищених маршрутів
- `GET /users`
- `GET /users/{id}`
- `POST /categories`
- `POST /records`
> Усі вони в групі під `AuthRequired`, отже потрібен `Authorization: Bearer <JWT>`.

---

## 8. Архітектура змін (коротко)

```
internal/
 ├─ http/
 │   ├─ handlers/
 │   │   ├─ auth.go            # Register/Login
 │   │   ├─ users.go, categories.go, records.go
 │   ├─ middleware/
 │   │   └─ auth.go            # AuthRequired: розбір Bearer, валідація JWT, userID у context
 │   └─ router.go              # публічні /auth/*, далі група під JWT
 ├─ security/
 │   └─ jwt.go                 # Sign(subject), Parse(token)
 └─ storage/
     ├─ storage.go             # +AuthUserStore interface
     └─ postgres/
         └─ users.go           # CreateWithPassword, FindAuth
migrations/
 ├─ 0001_init.sql
 ├─ 0002_seed.sql
 └─ 0003_add_password_hash.sql
```

**Важливо**: у `router.go` об’єкт `us` (типу `storage.UserStore`) приводиться до `storage.AuthUserStore`:
```go
aus, ok := us.(storage.AuthUserStore)
if !ok { panic("UserStore does not implement AuthUserStore") }
auth := handlers.AuthHandler{Users: aus}
```

---

## 9. Типові помилки та як фіксити
- `relation "categories" does not exist` → не накочені міграції. Застосуйте `0001/0002/0003`.
- `404 page not found` на `/auth/register` → сервіс зібраний без нових маршрутів. Зробіть **повну** пересборку:
  ```bash
  docker compose down
  docker compose build --no-cache sparrow-api
  docker compose up sparrow-api
  ```
- `401 Unauthorized` на закритих ендпоінтах → не переданий `Authorization: Bearer <JWT>` або токен протермінований/підпис інший.
- `pq: duplicate key value violates unique constraint` при повторній реєстрації → ім’я вже зайнято. Це очікувано.

---

## 10. Перевірочний чек‑лист для звіту (скріни)
1. `GET /health` → 200 OK без токена.
2. `POST /auth/register` → 201 Created.
3. `POST /auth/login` → відповідь з `access_token`.
4. `GET /users` без токена → 401.
5. `GET /users` з токеном → 200 OK.
6. `POST /categories` з токеном → 201 Created.
7. `POST /records` з токеном → 201 Created.

---

## Висновки
У рамках ЛР‑4 реалізовано повний цикл аутентифікації на базі JWT: реєстрація користувача з хешуванням пароля, видача токену при вході та захист основних ресурсів API через middleware. Додано міграції БД, налаштовано змінні середовища для секретів і TTL токену, підготовлено сценарії тестування (Postman/cURL). Отримане рішення забезпечує чіткий поділ публічних і приватних ендпоінтів та готове до розгортання у прод‑середовищі.

---

© KPI — Pliffdax – Sparrow API — 2025
