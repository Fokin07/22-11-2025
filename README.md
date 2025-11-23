# LinkChecker Service

Веб-сервис для проверки доступности интернет-ресурсов и генерации отчетов в формате PDF. Сервис поддерживает graceful shutdown/restart с сохранением состояния

# Возможности

**Проверка ссылок** - синхронная проверка доступности интернет-ресурсов

**Параллельная обработка** - одновременные проверки

**PDF отчеты** - генерация отчетов по ранее проверенным ссылкам

**Graceful shutdown** - корректное завершение с сохранением состояния

**Graceful restart** - перезапуск без потери данных

**Сохранение состояния** - восстановление данных после перезапуска


# Запуск
1. Создайте переменные окружения со значениями:     
    * AUTH_PORT=":8080"
2. Склонируйте этот репозиторий:
```
git clone https://github.com/Fokin07/22-11-2025.git
```
3. Соберите проект: `make build`. В консоли вы увидите:
```
go build -C ./cmd -o ./bin/app
```
4. Запустите сервис: `./cmd/bin/app`

# Использование
Проверка ссылок:
```
curl -X POST http://localhost:8080/check -H "Content-Type: application/json" -d '{"links": ["google.com", "yandex.ru", "invalid-domain.ru"]}'
```
```
curl -X POST http://localhost:8080/check -H "Content-Type: application/json" -d '{"links": ["habr.com", "invalid-domain.com","work-mate.ru"]}'
```

Генерация отчета:
```
curl -X GET http://localhost:8080/report -H "Content-Type: application/json" -d '{"links_list": [1, 2]}' --output report.pdf
```

# Паттерны

**1. Многослойная архитектура (Layered Architecture):**

    Presentation Layer (HTTP Handlers)

    Business Logic Layer (Services)

    Data Access Layer (Repository)

    Data Layer (Models)


**2. Dependency Injection (Внедрение зависимостей)**

**3. Graceful Shutdown/Restart Patterns**


# Подходы/Принципы
**1. SOLID**

**2. DRY**

**3. KISS**

**4. YAGNI**


# Endpoints
**POST /check**

Проверяет доступность ссылок

Request:
```
json
{
    "links": ["google.com", "invalid-domain.com"]
}
```

Response:
```
json
{
    "links": {
        "google.com": "available",
        "invalid-domain.com": "not available"
    },
    "links_num": 1
}
```
Request:
```
json
{
    "links": ["yandex.ru", "habr.com"]
}
```

Response:
```
json
{
    "links": {
        "google.com": "available",
        "habr.com": "available"
    },
    "links_num": 2
}
```

**GET /report**

Генерирует PDF отчет

Request:
```
json
{
    "links_list": [1, 2]
}
```

Response: 
```
PDF файл с отчётом
```