# Тестовое задание в Ozon


## Варианты запуска

### Вручную

Запуск in-memory с GraphQL playground а порту 8001
```bash
go run ./cmd -storage-type=false -debug=true -port=8001 
```

Для запуска с БД, указывается ```-storage-type=true```, а все параметры для подключения к БД указываются в .env и требуют перед запуском
```bash
source .env
```

### docker

Запуск in-memory с GraphQL playground а порту 8001
```bash
docker build . -t ozon-task
docker run -p 8001:8001 ozon-task ./main -storage-type=false -debug=true -port=8001 
```

### docker-compose

Запуск с БД с проброшенным портом 8080
```bash
docker-compose up -d
```
