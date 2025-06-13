# Tasks API

Простое HTTP-API для работы со списком задач (CRUD + упорядоченный вывод).

## 📦 Требования

- Go 1.24+  
- Git  


## ⚙️ Установка

1. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/ArtSEfimov/golang_education_task_3.git
   cd golang_education_task_3
   ```
2. Подтяните зависимости:
   ```bash
   go mod tidy
   ```
   Эта команда создаст/обновит файлы `go.mod` и `go.sum`, чтобы в проекте были все нужные пакеты.

## 🔧 Конфигурация

Сервис загружает следующие переменные окружения из файла `.env`:

```dotenv
# Порт сервера
PORT=8080
```

## 🚀 Запуск

1. Соберите и запустите сервер:
   ```bash
   go run main.go
   ```
   или
   ```bash
   go build -o bin/tasks-api main.go
   ./bin/tasks-api
   ```
2. По умолчанию сервер будет слушать порт из `PORT` (например, `http://localhost:8080`).

## 🔗 Эндпоинты

| Метод | Путь            | Описание                                            |
|-------|-----------------|-----------------------------------------------------|
| GET   | `/tasks`        | Получить все задачи. Параметр `ordered` (bool) — упорядоченный список. |
| GET   | `/tasks/{id}`   | Получить задачу по идентификатору.                  |
| POST  | `/tasks`        | Создать новую задачу.                               |
| PUT   | `/tasks/{id}`   | Обновить существующую задачу.                       |
| DELETE| `/tasks/{id}`   | Удалить задачу.                                     |

## ✏️ Примеры запросов

### 1. Получить все задачи (неупорядоченно)

**PowerShell**
```bash
Invoke-RestMethod `
  -Uri 'http://localhost:8080/tasks' `
  -Method GET
```

**curl.exe**
```bash
curl.exe http://localhost:8080/tasks
```

### 2. Получить все задачи в порядке добавления

```bash
curl http://localhost:8080/tasks?ordered=true
```

### 3. Получить задачу по ID

```bash
curl http://localhost:8080/tasks/42
```

### 4. Создать новую задачу

```bash
curl -X POST http://localhost:8080/tasks   -H "Content-Type: application/json"   -d '{"title":"Новая задача","description":"Описание задачи"}'
```

Тело запроса (JSON):
```json
{
  "title": "Новая задача",
  "description": "Дополнительное описание (необязательно)"
}
```

### 5. Обновить задачу

```bash
curl -X PUT http://localhost:8080/tasks/42   -H "Content-Type: application/json"   -d '{"title":"Обновлённый заголовок","description":"Обновлённое описание"}'
```

### 6. Удалить задачу

```bash
curl -X DELETE http://localhost:8080/tasks/42
```
