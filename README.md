# go-tasks-api
REST API for task management with async logging

## Запуск

```bash
git clone https://github.com/<your-username>/go-tasks-api.git
cd go-tasks-api
go run ./cmd/server
```

### Примеры запросов 
### Создать задачу
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Learn Go","description":"Study smth","status":"todo"}'
```

### Получить список:
```bash
curl http://localhost:8080/tasks
```
