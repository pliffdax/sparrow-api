# Лабораторна робота №1

## Мета роботи
Налаштувати середовище для подальшої розробки, створити базовий REST API з endpoint `/health`, контейнеризувати застосунок та задеплоїти його у хмару.

## Основні посилання
- Репозиторій проекту: [`sparrow-api`](https://github.com/pliffdax/sparrow-api)
- Деплой проекту: [`Render`](https://sparrow-api-l8pp.onrender.com/health)

## Стек технологій
- Мова програмування: Go (chi, net/http)
- Контейнеризація: Docker, Docker Compose
- Хостинг: Render
- Контроль версій: Git + Conventional Commits

## Хід виконання
1. Ініціалізовано Go-проєкт, підключено бібліотеку `chi`.
2. Реалізовано endpoint `GET /health`, який повертає JSON зі статусом, часом та версією.
3. Створено `Dockerfile` (builder + runtime), зібрано образ і перевірено запуск.
4. Створено `docker-compose.yml` для зручного керування сервісом.
5. Виконано локальне тестування через `curl`.
6. Налаштовано деплой на Render. Сервіс доступний за адресою:  
   👉 [https://sparrow-api-l8pp.onrender.com/health](https://sparrow-api-l8pp.onrender.com/health)

## Перевірка
```bash
curl http://localhost:8080/health
# -> {"status":200,"time":"2025-10-01T12:56:19Z","version":"0.1.0"}
```

## Висновки

У ході виконання роботи було створено базовий веб-сервіс на Go, налаштовано контейнеризацію через Docker та Docker Compose, а також успішно здійснено деплой на Render. Отримані навички можна застосувати для швидкого створення та публікації backend-сервісів.

---

© KPI — Pliffdax – Sparrow API — 2025