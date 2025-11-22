# LinkChecker Service

Веб-сервис для проверки доступности интернет-ресурсов и генерации отчетов в формате PDF. Сервис поддерживает graceful shutdown/restart с сохранением состояния

# Возможности

**Проверка ссылок** - синхронная проверка доступности интернет-ресурсов

**Параллельная обработка** - до 10 одновременных проверок

**PDF отчеты** - генерация отчетов по ранее проверенным ссылкам

**Graceful shutdown** - корректное завершение с сохранением состояния

**Graceful restart** - перезапуск без потери данных

**Сохранение состояния** - восстановление данных после перезапуска

# Паттерны

**1. Многослойная архитектура (Layered Architecture):**

    Presentation Layer (HTTP Handlers)

    Business Logic Layer (Services)

    Data Access Layer (Repository)

    Data Layer (Models)


**2. Dependency Injection (Внедрение зависимостей)**

**3. Graceful Shutdown/Restart Patterns**

