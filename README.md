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
| GET   | `/tasks`        | Получить все задачи.                                |
| GET   | `/tasks/{id}`   | Получить задачу по идентификатору.                  |
| POST  | `/tasks`        | Создать новую задачу.                               |
| PUT   | `/tasks/{id}`   | Обновить существующую задачу.                       |
| DELETE| `/tasks/{id}`   | Удалить задачу.                                     |


## 🛠 Опции запросов

| Параметр   | Описание                                | Пример запроса                                                                     |
|------------|-----------------------------------------|------------------------------------------------------------------------------------|
| `ordered`  | Упорядоченный список                   | `GET /tasks?ordered=true` — вернуть задачи в порядке создания                       |
| `created`  | Включить задачи в статусе **Created**  | `GET /tasks?created=true` — показать только только что созданные задачи            |
| `running`  | Включить задачи в статусе **Running**  | `GET /tasks?running=true` — показать только задачи, находящиеся в процессе выполнения |
| `completed`| Включить задачи в статусе **Completed**| `GET /tasks?completed=true` — показать только успешно завершённые задачи           |
| `failed`   | Включить задачи в статусе **Failed**   | `GET /tasks?failed=true` — показать только задачи, завершившиеся с ошибкой         |

## Примеры комбинированных запросов

### 1. Упорядоченный список только завершённых задач
GET /tasks?ordered=true&completed=true

### 2. Список запущенных задач
GET /tasks?running=true

### 3. Только недавно созданные задачи, не упорядоченные
GET /tasks?created=true&ordered=false

### 4. Все задачи, кроме завершённых и неудачных
GET /tasks?completed=false&failed=false

### 5. Упорядоченный список запущенных и только что созданных задач
GET /tasks?ordered=true&running=true&created=true

### 6. Показать только неупорядоченные и не созданные задачи
GET /tasks?ordered=false&created=false

### 7. Список задач по любым параметрам (все включены)
GET /tasks?ordered=true&created=true&running=true&completed=true&failed=true

### 8. Без фильтров (вернуть всё)
GET /tasks


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

**PowerShell**
```bash
Invoke-RestMethod `
  -Uri 'http://localhost:8080/tasks?ordered=true' `
  -Method GET
```

**curl.exe**
```bash
curl.exe "http://localhost:8080/tasks?ordered=true"
```

### 3. Получить задачу по ID

**PowerShell**
```bash
$id = 42
Invoke-RestMethod `
  -Uri "http://localhost:8080/tasks/$id" `
  -Method GET
```

**curl.exe**
```bash
$id=42
curl.exe http://localhost:8080/tasks/$id
```

### 4. Создать новую задачу

**PowerShell**
```bash
$body = @{
  title       = 'Новая задача'
  description = 'Описание задачи'
}

Invoke-RestMethod `
  -Uri 'http://localhost:8080/tasks' `
  -Method POST `
  -ContentType 'application/json' `
  -Body (ConvertTo-Json $body)
```

**curl.exe**
```bash
curl.exe -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Новая задача","description":"Описание задачи"}'
```

Тело запроса (JSON):
```bash
{
  title       = 'Новая задача'
  description = 'Описание задачи' (необязательное поле)
}
```

### 5. Обновить задачу

**PowerShell**
```bash
$id = 42

$body = @{
  title       = 'Обновлённый заголовок'
  description = 'Обновлённое описание'
}

Invoke-RestMethod `
  -Uri "http://localhost:8080/tasks/$id" `
  -Method PUT `
  -ContentType 'application/json' `
  -Body (ConvertTo-Json $body)
```

**curl.exe**
```bash
curl.exe -X PUT http://localhost:8080/tasks/42 \
  -H "Content-Type: application/json" \
  -d '{"title":"Обновлённый заголовок","description":"Обновлённое описание"}'
```


### 6. Удалить задачу

**PowerShell**
```bash
$id = 42

Invoke-RestMethod `
  -Uri "http://localhost:8080/tasks/$id" `
  -Method DELETE
```

**curl.exe**
```bash
$id=42
curl.exe -X DELETE http://localhost:8080/tasks/$id
```
