# Subscription Service API

REST-сервис для агрегации данных об онлайн-подписках пользователей. Выполнено в рамках тестового задания на позицию Junior Golang Developer (Effective Mobile).

## Стек технологий
* **Язык:** Go
* **Фреймворк:** Echo
* **База данных:** PostgreSQL + GORM
* **Миграции:** Goose
* **Логирование:** `log/slog`
* **Инфраструктура:** Docker, Docker Compose, Makefile

## Запуск сервиса

Убедитесь, что у вас установлены Docker и Docker Compose.

1. Склонируйте репозиторий:
```bash
git clone <ссылка_на_твой_репозиторий>
cd <название_папки>
```

2. Запустите проект с помощью `make docker-build`. (поднимется БД, автоматически накатятся миграции, запустится приложение) Сервис будет доступен на порту 1323.

## Документация (Swagger)

Документация Swagger со всеми схемами запросов и ответов доступна по адресу:

-> **http://localhost:1323/swagger/index.html**

## Примеры использования (cURL)

<details>
<summary><b>1. Создать новую подписку</b></summary>

```bash
curl -X POST http://localhost:1323/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025"
}'
```
</details>

<details>
<summary><b>2. Получить список всех подписок (List)</b></summary>

```bash
curl -X GET http://localhost:1323/subscriptions
```
</details>

<details>
<summary><b>3. Подсчет суммарной стоимости</b></summary>

Подсчет стоимости всех подписок пользователя за выбранный период.

```bash
curl -X GET "http://localhost:1323/subscriptions/sum?user_id=94f639ba-6d37-484a-b26a-d1260cc99ef8&start_date=07-2025"
```
</details>

*(Остальные методы CRUDL детально описаны в Swagger).*

## Разработка и тестирование (Makefile)

Для удобства локальной разработки предусмотрен `Makefile`:

* `make docker-build` — Полный запуск проекта в Docker с пересборкой образов.
* `make docker-up` — Обычный запуск контейнеров (без билда).
* `make docker-down` — Остановка контейнеров и удаление volumes (очистка БД).
* `make run-local` — Локальный запуск приложения (требует запущенной БД на порту 5433).
* `make test` — Запуск всех тестов проекта.
