# Лабораторна робота №3 — Валідація, обробка помилок, ORM

## Мета роботи
Покращити застосунок Sparrow API шляхом додавання валідації вхідних даних, централізованої обробки помилок та переходу на ORM (GORM) з базою даних PostgreSQL.

---

## 1️⃣ Підготовка середовища

### Docker Compose
У проєкт додано сервіс бази даних PostgreSQL:

```yaml
services:
  db:
    image: postgres:16
    environment:
      POSTGRES_DB: ${POSTGRES_DB}
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
    volumes:
      - pgdata:/var/lib/postgresql/data
      - ./migrations:/migrations
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER} -d ${POSTGRES_DB}"]
      interval: 3s
      timeout: 3s
      retries: 20
  sparrow-api:
    build: .
    environment:
      PORT: 8080
      DB_DSN: ${DB_DSN}
    depends_on:
      db:
        condition: service_healthy
    ports:
      - "8080:8080"

volumes:
  pgdata:
```

### Файл `.env`
```
POSTGRES_DB=sparrow
POSTGRES_USER=sparrow
POSTGRES_PASSWORD=secret
DB_DSN=postgres://sparrow:secret@db:5432/sparrow?sslmode=disable
PORT=8080
```

---

## 2️⃣ Міграції та сідінг

### `0001_migration.sql`
Створює таблиці `users`, `categories`, `records`:

```sql
CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS categories (
  id BIGSERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  is_global BOOLEAN NOT NULL DEFAULT FALSE,
  owner_id BIGINT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS records (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  category_id BIGINT NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
  amount NUMERIC(14,2) NOT NULL CHECK (amount > 0),
  happened_at TIMESTAMPTZ NOT NULL,
  note TEXT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

### `0002_seed.sql`
Створює базові дані (користувачів і глобальні категорії):
```sql
INSERT INTO users (name)
SELECT 'Demo' WHERE NOT EXISTS (SELECT 1 FROM users WHERE name='Demo');
INSERT INTO users (name)
SELECT 'Other' WHERE NOT EXISTS (SELECT 1 FROM users WHERE name='Other');

INSERT INTO categories (title, is_global)
SELECT 'Food', TRUE WHERE NOT EXISTS (SELECT 1 FROM categories WHERE title='Food' AND is_global=TRUE);
INSERT INTO categories (title, is_global)
SELECT 'Transport', TRUE WHERE NOT EXISTS (SELECT 1 FROM categories WHERE title='Transport' AND is_global=TRUE);
INSERT INTO categories (title, is_global)
SELECT 'Health', TRUE WHERE NOT EXISTS (SELECT 1 FROM categories WHERE title='Health' AND is_global=TRUE);
```

---

## 3️⃣ Варіант завдання

**Група:** IO-35  
**Варіант:** 35 % 3 = 2 → *Користувацькі категорії витрат*

Реалізовано:
- Глобальні категорії (`is_global = TRUE, owner_id = NULL`) — видимі всім користувачам;
- Користувацькі категорії (`is_global = FALSE, owner_id = user_id`) — доступні лише власнику;
- Обмеження доступу при видаленні чужих категорій (`403 Forbidden`);
- Перевірка при створенні витрат — користувач може використати лише свої категорії.

---

## 4️⃣ Тестування API (Postman)

### Environment `Local`
```
baseUrl = http://localhost:8080
userId  = 1
```

### Запити

| № | Метод | Endpoint | Опис | Очікувано |
|---|--------|-----------|------|------------|
| 1 | GET | `/categories` | Отримати всі категорії (глобальні + власні) | 200 OK |
| 2 | POST | `/categories` | Створити власну категорію | 201 Created |
| 3 | DELETE | `/categories/{id}` | Видалити свою категорію | 204 No Content |
| 4 | DELETE | `/categories/{id}` з іншим X-User-ID | Спроба видалити чужу категорію | 403 Forbidden |
| 5 | POST | `/records` | Створити витрату зі своєю категорією | 201 Created |
| 6 | POST | `/records` з чужою категорією | Заборонена дія | 403 Forbidden |

---

## 5️⃣ Команди запуску

```bash
# Підняти середовище
docker compose up -d --build

# Застосувати міграції та сідінг
docker compose exec -T db psql -U $POSTGRES_USER -d $POSTGRES_DB -f /migrations/0001_migration.sql
docker compose exec -T db psql -U $POSTGRES_USER -d $POSTGRES_DB -f /migrations/0002_seed.sql

# Перевірити дані
docker compose exec -T db psql -U $POSTGRES_USER -d $POSTGRES_DB -c "SELECT * FROM users;"
docker compose exec -T db psql -U $POSTGRES_USER -d $POSTGRES_DB -c "SELECT * FROM categories;"
```

---

## 6️⃣ Результати тестів

- ✅ Глобальні категорії видимі всім користувачам.  
- ✅ Користувацькі категорії видимі лише власнику.  
- ✅ Чужі категорії видалити не можна (403).  
- ✅ Неможливо створити витрату з чужою категорією (403).

---

## 7️⃣ Теги версій

| ЛР | Версія | Коментар |
|----|---------|-----------|
| Лабораторна 1 | v1.0.0 | Базові endpoints |
| Лабораторна 2 | v2.0.0 | Postman, Docker, документація |
| Лабораторна 3 | v3.0.0 | ORM, валідація, обробка помилок, користувацькі категорії |

## Висновки

У ході виконання лабораторної роботи було інтегровано підтримку бази даних PostgreSQL за допомогою ORM-бібліотеки GORM, що дозволило відмовитись від in-memory сховищ і забезпечити збереження даних між запусками сервера. Реалізовано централізовану обробку помилок та базову валідацію вхідних даних. Для тестування API створено Postman-колекцію, яка перевіряє коректність CRUD-операцій та логіку доступу до користувацьких категорій. Згідно з варіантом 2, додано підтримку глобальних і користувацьких категорій витрат із відповідними обмеженнями доступу. Розроблене рішення є стабільною основою для подальшого розширення функціоналу застосунку.

---

© KPI — Pliffdax – Sparrow API — 2025
