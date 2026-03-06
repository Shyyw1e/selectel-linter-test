# selectel-linter-test

Кастомный линтер для Go, который валидирует лог-сообщения в вызовах `slog` и `zap`.  
Реализован на базе `go/analysis`, с точкой интеграции для `golangci-lint`.

## Что проверяет линтер

Правила проверки лог-сообщений:

1. Сообщение должно начинаться со строчной буквы.
2. Сообщение должно быть на английском языке (кириллица запрещена).
3. Сообщение не должно содержать спецсимволы и emoji.
4. Сообщение не должно содержать потенциально чувствительные данные.

Поддерживаемые логгеры:
- `log/slog`
- `go.uber.org/zap`

## Структура проекта

```text
internal/
  analyser/   # runtime анализатора + AST parsing/extraction
  config/     # конфигурация, дефолты, валидация
  rules/      # контракты и реализации правил
plugin/       # экспорт анализатора для интеграции
```

## Требования

- Go `1.24+`
- `golangci-lint` (для интеграционных проверок)

## Быстрый старт

```bash
go test ./...
```

## Реализованные правила

- `lowercase-start`
- `english-only`
- `no-special-chars`
- `no-sensitive-data`

Проверка чувствительных данных работает через:
- список ключевых слов (`SensitiveKeywords`)
- регулярные выражения (`SensitivePatterns`)

## Конфигурация

Основная структура (`internal/config/config.go`):

- `EnabledRules map[string]bool`
- `SensitiveKeywords []string`
- `SensitivePatterns []string`

Валидация конфигурации включает:
- проверку неизвестных `rule id`
- нормализацию списка ключевых слов
- проверку корректности regex-паттернов

## Пайплайн анализатора

1. Обход AST (`ast.Inspect`) и поиск `CallExpr`.
2. Распознавание поддерживаемых лог-вызовов (`slog`/`zap`) через type info.
3. Извлечение текста сообщения из:
   - строкового литерала
   - конкатенации литералов (`"a" + "b"`)
4. Применение включенных правил.
5. Репорт диагностик через `pass.Report(...)`.
6. Для правила lowercase добавляется `SuggestedFix` (для простых строковых литералов).

## Тестирование

Покрытие тестами:
- unit-тесты для каждого правила
- тесты extractor/parser
- интеграционный тест через `analysistest` (`internal/analyser/testdata`)

Запуск:

```bash
go test ./...
```

## CI

GitHub Actions workflow (`.github/workflows/ci.yml`) запускается при:
- push в `main/master`
- pull request

Шаги pipeline:
1. checkout
2. setup Go `1.24.x`
3. `go mod download`
4. проверка форматирования (`gofmt`)
5. `go test ./...`
6. `golangci-lint`

## Ограничения

- Извлечение текста пока поддерживает только литералы и их конкатенацию.
- Динамические сообщения (`fmt.Sprintf`, переменные, сложные выражения) сознательно пропускаются.
- По дизайну поддерживаются только `slog` и `zap`.

## Интеграция

Точка входа плагина: `plugin/plugin.go`.  
Анализатор создается в `internal/analyser`, конфигурация берется из `internal/config`.

## Примеры

Нарушения:

```go
slog.Info("Starting server")
slog.Info("запуск сервера")
slog.Warn("connection failed!!!")
slog.Debug("token validated")
```

Корректные примеры:

```go
slog.Info("starting server")
slog.Warn("connection failed")
slog.Info("request completed")
```

## Запуск через golangci-lint (Module Plugin System)

### 1. Подготовить конфиги

В корне репозитория должны быть файлы:

- `.custom-gcl.yml`
- `.golangci.yml`

### 2. Собрать кастомный бинарь golangci-lint

```bash
golangci-lint custom
```

После команды появится бинарь `./custom-golangci-lint`.

### 3. Запустить проверки

```bash
./custom-golangci-lint run ./...
```

### 4. Обновить бинарь после изменений

Если изменился код линтера или конфиги плагина, пересобери бинарь:

```bash
golangci-lint custom
```
```

И если хочешь, вставь в `README` коротко про конфиг:

```md
Линтер зарегистрирован под именем `loglint` и включается в `.golangci.yml` через:

```yaml
linters:
  enable:
    - loglint
```