# booker-admin-gateway

REST API Gateway для системы бронирования столов.

## Описание

Admin Gateway предоставляет единый REST API для управления заведениями и бронированиями. Проксирует запросы к внутренним gRPC сервисам (venue-svc, booking-svc).

## Запуск

```bash
go run cmd/admin-gateway/main.go
```

## Docker

```bash
docker build -t ghcr.io/bookingcontrol/booker-admin-gateway:v1.0.0 .
```

## Зависимости

- `github.com/bookingcontrol/booker-contracts-go/v1` - Protobuf контракты
