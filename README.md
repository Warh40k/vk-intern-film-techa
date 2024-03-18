# API-сервис Фильмотека

Для запуска собрать docker-образ в корне проекта:
```go
docker build --tag=filmotecka:latest .
docker compose up -d
```

Приложение будет доступно на 8080 порту. Документация Swagger - на порту 8000 и в каталоге docs.