# Admin Gateway (admin-gateway)

## 📋 Обзор

**Admin Gateway** - это API Gateway для административной панели системы управления бронированиями. Сервис предоставляет HTTP REST API для фронтенда и проксирует запросы к gRPC backend-сервисам (venue-svc и booking-svc).

**Технологии:**
- Go 1.23
- Echo Framework (HTTP router)
- gRPC clients (для backend-сервисов)
- Redis (аутентификация, rate limiting, сессии)
- Jaeger (distributed tracing)
- Prometheus (метрики)

**Порт:**
- HTTP API: 8080 (внешний 18080)

---

## 🏗️ Архитектура

### Слои приложения

Сервис построен по **Clean Architecture** с Hexagonal (Ports & Adapters) подходом:

```
admin-gateway/
├── cmd/admin-gateway/
│   └── main.go                      # Точка входа
├── internal/
│   ├── adapter/                     # Адаптеры (внешние интерфейсы)
│   │   ├── http/                   # HTTP handlers (входящие)
│   │   │   ├── auth_handler.go
│   │   │   ├── booking_handler.go
│   │   │   ├── venue_handler.go
│   │   │   ├── router.go
│   │   │   └── middleware/         # HTTP middleware
│   │   ├── grpc/                   # gRPC clients (исходящие)
│   │   │   ├── booking_repo.go
│   │   │   └── venue_repo.go
│   │   └── redis/                  # Redis client (исходящие)
│   │       └── auth_repo.go
│   ├── domain/                      # Доменные интерфейсы (порты)
│   │   ├── auth/
│   │   │   └── repository.go      # Auth repository interface
│   │   ├── booking/
│   │   │   └── repository.go      # Booking repository interface
│   │   └── venue/
│   │       └── repository.go      # Venue repository interface
│   ├── usecase/                     # Бизнес-логика (use cases)
│   │   ├── auth/
│   │   │   └── service.go
│   │   ├── booking/
│   │   │   └── service.go
│   │   └── venue/
│   │       └── service.go
│   ├── infrastructure/              # Инфраструктурные компоненты
│   │   ├── redis/
│   │   │   └── client.go
│   │   ├── tracing/
│   │   │   └── tracing.go
│   │   └── metrics/
│   │       └── metrics.go
│   └── config/
│       └── config.go
└── Makefile
```

### Принципы

1. **Hexagonal Architecture** - чистое разделение бизнес-логики и адаптеров
2. **Ports & Adapters** - domain определяет интерфейсы (порты), adapters их реализуют
3. **Dependency Inversion** - domain не зависит от infrastructure
4. **Single Responsibility** - каждый handler отвечает за один ресурс

### Слои взаимодействия

```
HTTP Request
     ↓
[HTTP Handler]  ← adapter/http/
     ↓
[Use Case]      ← usecase/ (business logic)
     ↓
[Repository]    ← domain/ (interface/port)
     ↓
[Adapter]       ← adapter/grpc/ или adapter/redis/ (реализация)
     ↓
External Service (gRPC, Redis)
```

---

## 🔄 Межсервисное взаимодействие

### Входящие запросы (HTTP REST API)

Admin Gateway предоставляет REST API:

#### Authentication
```
POST   /api/auth/register       # Регистрация пользователя
POST   /api/auth/login          # Вход (получение токена)
POST   /api/auth/refresh        # Обновление токена
```

#### Venues
```
POST   /api/venues              # Создать заведение
GET    /api/venues              # Список заведений
GET    /api/venues/:id          # Получить заведение
PUT    /api/venues/:id          # Обновить заведение
DELETE /api/venues/:id          # Удалить заведение
```

#### Rooms
```
POST   /api/venues/:venueId/rooms           # Создать зал
GET    /api/venues/:venueId/rooms           # Список залов
GET    /api/rooms/:id                       # Получить зал
PUT    /api/rooms/:id                       # Обновить зал
DELETE /api/rooms/:id                       # Удалить зал
```

#### Tables
```
POST   /api/rooms/:roomId/tables            # Создать стол
GET    /api/rooms/:roomId/tables            # Список столов в зале
GET    /api/venues/:venueId/tables          # Все столы в заведении
GET    /api/tables/:id                      # Получить стол
PUT    /api/tables/:id                      # Обновить стол
DELETE /api/tables/:id                      # Удалить стол
```

#### Bookings
```
POST   /api/bookings                        # Создать бронирование
GET    /api/bookings                        # Список бронирований
GET    /api/bookings/:id                    # Получить бронирование
POST   /api/bookings/:id/confirm            # Подтвердить бронирование
POST   /api/bookings/:id/cancel             # Отменить бронирование
POST   /api/bookings/:id/seated             # Отметить посадку
POST   /api/bookings/:id/finished           # Завершить бронирование
POST   /api/bookings/:id/no-show            # Отметить no-show
```

#### Availability
```
POST   /api/venues/:venueId/check-availability  # Проверить доступность столов
```

#### Static files
```
GET    /                        # Serve index.html
GET    /assets/*                # Serve static assets (JS, CSS)
```

### Исходящие запросы

#### gRPC клиенты

**Venue Service:**
```go
type VenueRepo struct {
    client venuepb.VenueServiceClient
}

func (r *VenueRepo) CreateVenue(ctx context.Context, in CreateVenueInput) (VenueView, error) {
    resp, err := r.client.CreateVenue(ctx, &venuepb.CreateVenueRequest{
        Name:     in.Name,
        Timezone: in.Timezone,
        Address:  in.Address,
        Phone:    in.Phone,
        Email:    in.Email,
    })
    return toVenueView(resp), err
}
```

**Booking Service:**
```go
type BookingRepo struct {
    client bookingpb.BookingServiceClient
}

func (r *BookingRepo) CreateBooking(ctx context.Context, in CreateBookingInput) (BookingView, error) {
    resp, err := r.client.CreateBooking(ctx, &bookingpb.CreateBookingRequest{
        VenueId:       in.VenueID,
        Table:         &commonpb.TableRef{...},
        Slot:          &commonpb.Slot{...},
        PartySize:     in.PartySize,
        CustomerName:  in.CustomerName,
        CustomerPhone: in.CustomerPhone,
        Comment:       in.Comment,
        AdminId:       in.AdminID,
    })
    return toBookingView(resp), err
}
```

#### Redis

Используется для:
1. **Аутентификация** - хранение пользователей (MVP, в production должна быть БД)
2. **Сессии** - хранение активных токенов
3. **Rate Limiting** - ограничение запросов

```go
type AuthRepo struct {
    redis *redis.Client
}

func (r *AuthRepo) CreateUser(ctx context.Context, username string, data map[string]interface{}) error {
    key := fmt.Sprintf("user:%s", username)
    jsonData, _ := json.Marshal(data)
    return r.redis.Set(ctx, key, jsonData, 0).Err()
}

func (r *AuthRepo) UserExists(ctx context.Context, username string) (bool, error) {
    key := fmt.Sprintf("user:%s", username)
    exists, err := r.redis.Exists(ctx, key).Result()
    return exists > 0, err
}
```

---

## 🔒 Аутентификация и авторизация

### Текущая реализация (MVP)

**⚠️ WARNING: Текущая реализация НЕ для production!**

Для MVP используется упрощенная аутентификация:

1. **Пользователи хранятся в Redis** (должна быть БД)
2. **Пароли НЕ хэшируются** (должен быть bcrypt/argon2)
3. **Токены простые строки** (должны быть JWT)
4. **Нет refresh token механизма** (должен быть)

```go
func (s *Service) Login(ctx context.Context, in LoginInput) (LoginView, error) {
    // Проверка пользователя
    exists, _ := s.repo.UserExists(ctx, in.Username)
    if !exists {
        return LoginView{}, errors.New("invalid credentials")
    }
    
    // Проверка пароля (НЕ хэшированный!)
    storedPassword, _ := s.repo.GetUserPassword(ctx, in.Username)
    if storedPassword != in.Password {
        return LoginView{}, errors.New("invalid credentials")
    }
    
    // Генерация токена (простая строка!)
    token := "token-" + in.Username
    refreshToken := "refresh-" + in.Username
    
    return LoginView{
        AccessToken:  token,
        RefreshToken: refreshToken,
    }, nil
}
```

### Middleware

```go
func (m *Middleware) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        // Извлечение токена из заголовка
        authHeader := c.Request().Header.Get("Authorization")
        if authHeader == "" {
            return c.JSON(401, map[string]string{"error": "missing authorization header"})
        }
        
        // Формат: "Bearer {token}"
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            return c.JSON(401, map[string]string{"error": "invalid authorization header"})
        }
        
        token := parts[1]
        
        // TODO: Валидация JWT
        adminID := "admin-1"  // Хардкод для MVP
        
        // Сохранение в контекст
        c.Set("admin_id", adminID)
        c.Set("token", token)
        
        return next(c)
    }
}
```

### Rate Limiting

```go
func (m *Middleware) RateLimitMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            adminID := c.Get("admin_id")
            if adminID == nil {
                return next(c)  // Не аутентифицирован - пропускаем
            }
            
            // Ключ в Redis
            key := "rl:" + adminID.(string)
            limit := 100  // Запросов в минуту
            
            // Инкремент счетчика
            count, _ := m.redisClient.Incr(c.Request().Context(), key)
            
            // Установка TTL при первом запросе
            if count == 1 {
                m.redisClient.Expire(c.Request().Context(), key, time.Minute)
            }
            
            // Проверка лимита
            if count > int64(limit) {
                return c.JSON(429, map[string]string{"error": "rate limit exceeded"})
            }
            
            return next(c)
        }
    }
}
```

---

## 🗄️ Работа с Redis

### Структура данных

#### Пользователи

```
Ключ:  user:{username}
Значение: JSON с данными пользователя
TTL: нет

Пример:
user:admin = {"username":"admin","password":"admin123","email":"admin@example.com"}
```

#### Rate Limiting

```
Ключ:  rl:{admin_id}
Значение: счетчик запросов
TTL: 1 минута

Пример:
rl:admin-1 = 42 (TTL: 18 seconds)
```

### Операции

```go
type Client struct {
    *redis.Client
}

// Пользователи
func (c *Client) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
func (c *Client) Get(ctx context.Context, key string) (string, error)
func (c *Client) Exists(ctx context.Context, key string) (int64, error)

// Rate limiting
func (c *Client) Incr(ctx context.Context, key string) (int64, error)
func (c *Client) Expire(ctx context.Context, key string, ttl time.Duration) error
```

---

## 💼 Бизнес-логика

### Use Cases

Use case слой содержит бизнес-логику и оркестрацию вызовов к репозиториям.

#### Auth Use Case

```go
type Service struct {
    repo dom.Repository  // domain/auth/repository.go (interface)
}

func (s *Service) Register(ctx context.Context, in CreateInput) (RegisterView, error) {
    // Валидация
    if in.Username == "" || in.Password == "" {
        return RegisterView{}, errors.New("username and password are required")
    }
    
    // Проверка существования
    exists, _ := s.repo.UserExists(ctx, in.Username)
    if exists {
        return RegisterView{}, errors.New("username already exists")
    }
    
    // Создание пользователя
    userData := map[string]interface{}{
        "username": in.Username,
        "password": in.Password,  // TODO: Hash
        "email":    in.Email,
    }
    
    err := s.repo.CreateUser(ctx, in.Username, userData)
    if err != nil {
        return RegisterView{}, errors.New("failed to create user")
    }
    
    return RegisterView{
        Username: in.Username,
        Message:  "User registered successfully",
    }, nil
}
```

#### Venue Use Case

```go
type Service struct {
    repo dom.Repository  // domain/venue/repository.go (interface)
}

func (s *Service) CreateVenue(ctx context.Context, in CreateVenueInput) (VenueView, error) {
    // Простая проксация к venue-svc через repo
    return s.repo.CreateVenue(ctx, in)
}

func (s *Service) ListVenues(ctx context.Context, limit, offset int32) (ListVenuesView, error) {
    return s.repo.ListVenues(ctx, limit, offset)
}
```

Use case для venue и booking в основном просто проксирует запросы к gRPC сервисам, добавляя минимальную логику (логирование, валидацию).

---

## 📝 Логирование

### Библиотека

Используется **zerolog** для структурированного логирования.

### Примеры логов

```go
// Info - успешные операции
log.Info().
    Str("username", in.Username).
    Msg("User registered")

log.Info().
    Str("path", c.Path()).
    Str("method", c.Request().Method).
    Str("admin_id", adminID).
    Msg("AuthMiddleware: request authorized")

// Warning - потенциальные проблемы
log.Warn().
    Str("path", c.Path()).
    Str("method", c.Request().Method).
    Msg("AuthMiddleware: missing authorization header")

// Error - ошибки операций
log.Error().
    Err(err).
    Msg("Failed to check user existence")

log.Error().
    Err(err).
    Msg("Rate limit check failed")
```

### Middleware логирование

Echo framework автоматически логирует все HTTP запросы:

```go
e.Use(middleware.Logger())
```

Формат:
```
{"time":"2024-01-01T12:00:00Z","level":"info","method":"POST","uri":"/api/bookings","status":200,"latency":45123456}
```

---

## 📊 Метрики

### Prometheus Metrics

Сервис экспортирует метрики на `/metrics`.

#### HTTP метрики

```go
// Количество HTTP запросов
http_requests_total{method="POST", path="/api/bookings", status="200", service="admin-gateway"}

// Latency запросов
http_request_duration_seconds{method="POST", path="/api/bookings", status="200", service="admin-gateway"}
```

#### gRPC клиент метрики

```go
// Количество исходящих gRPC запросов
grpc_client_requests_total{method="CreateBooking", status="ok", service="admin-gateway"}

// Latency
grpc_client_request_duration_seconds{method="CreateBooking", status="ok", service="admin-gateway"}
```

#### Redis метрики

```go
redis_operations_total{operation="get", service="admin-gateway"}
redis_operation_duration_seconds{operation="incr", service="admin-gateway"}
```

### Middleware для метрик

```go
func MetricsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        start := time.Now()
        
        err := next(c)
        
        duration := time.Since(start).Seconds()
        status := c.Response().Status
        
        HTTPRequestsTotal.WithLabelValues(
            c.Request().Method,
            c.Path(),
            strconv.Itoa(status),
            "admin-gateway",
        ).Inc()
        
        HTTPRequestDuration.WithLabelValues(
            c.Request().Method,
            c.Path(),
            strconv.Itoa(status),
            "admin-gateway",
        ).Observe(duration)
        
        return err
    }
}
```

---

## 🔍 Трейсинг

### OpenTelemetry + Jaeger

```go
shutdown, err := tracing.InitTracer("admin-gateway", cfg.JaegerEndpoint)
if err != nil {
    log.Fatal().Err(err).Msg("Failed to initialize tracer")
}
defer shutdown()
```

### Propagation

Trace context автоматически передается в gRPC запросы:

```
HTTP Request (trace_id=abc123)
     ↓
admin-gateway (span: HandleCreateBooking)
     ↓ gRPC call (propagate trace_id)
booking-svc (span: CreateBooking, parent=HandleCreateBooking)
     ↓ gRPC call
venue-svc (span: CheckAvailability, parent=CreateBooking)
```

Jaeger UI показывает полную цепочку запросов across services.

---

## 🧪 Тестирование

### Unit Tests

```
internal/
├── adapter/
│   ├── http/
│   │   ├── auth_handler_test.go
│   │   ├── booking_handler_test.go
│   │   └── venue_handler_test.go
│   └── grpc/
│       └── booking_repo_test.go
└── usecase/
    └── auth/
        └── service_test.go
```

#### Запуск тестов

```bash
go test ./...
go test -cover ./...
go test ./internal/usecase/auth -v
```

### Integration Tests

```go
func TestAuthHandler_Integration(t *testing.T) {
    // Setup Echo + handlers
    e := echo.New()
    handler := auth_handler.New(authService)
    e.POST("/api/auth/register", handler.Register)
    
    // Test request
    req := httptest.NewRequest(http.MethodPost, "/api/auth/register", body)
    rec := httptest.NewRecorder()
    
    e.ServeHTTP(rec, req)
    
    assert.Equal(t, http.StatusOK, rec.Code)
}
```

### E2E Tests

Полные тесты через HTTP API:

```bash
# Регистрация
curl -X POST http://localhost:18080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test123","email":"test@example.com"}'

# Вход
curl -X POST http://localhost:18080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test123"}'

# Создание venue
curl -X POST http://localhost:18080/api/venues \
  -H "Authorization: Bearer token-test" \
  -H "Content-Type: application/json" \
  -d '{"name":"My Restaurant","timezone":"Europe/Moscow","address":"..."}'
```

---

## ⚙️ Конфигурация

### Переменные окружения

```bash
# Server
ENV=development                       # development/production
PORT=8080                            # HTTP порт

# gRPC Services
GRPC_VENUE_ADDR=venue-svc:50051
GRPC_BOOKING_ADDR=booking-svc:50052

# Redis
REDIS_ADDR=redis-master:6379
REDIS_PASSWORD=redis_pass

# Auth
JWT_SECRET=your-secret-key-change-in-production

# Tracing
JAEGER_ENDPOINT=http://jaeger:14268/api/traces
```

### Config struct

```go
type Config struct {
    Env             string
    Port            int
    GRPCVenueAddr   string
    GRPCBookingAddr string
    RedisAddr       string
    RedisPassword   string
    JWTSecret       string
    JaegerEndpoint  string
}

func Load() *Config {
    return &Config{
        Env:             getEnv("ENV", "development"),
        Port:            getEnvInt("PORT", 8080),
        GRPCVenueAddr:   getEnv("GRPC_VENUE_ADDR", "localhost:50051"),
        GRPCBookingAddr: getEnv("GRPC_BOOKING_ADDR", "localhost:50052"),
        RedisAddr:       getEnv("REDIS_ADDR", "localhost:6379"),
        RedisPassword:   getEnv("REDIS_PASSWORD", ""),
        JWTSecret:       getEnv("JWT_SECRET", "secret"),
        JaegerEndpoint:  getEnv("JAEGER_ENDPOINT", ""),
    }
}
```

---

## 🚀 Запуск

### Локальная разработка

```bash
# 1. Запустить инфраструктуру + backend сервисы
cd ../infra
docker compose --profile infra-min --profile apps up -d

# 2. Установить зависимости
go mod download

# 3. Запустить gateway
go run cmd/admin-gateway/main.go

# API доступен на http://localhost:8080
```

### С фронтендом

```bash
# 1. Собрать фронтенд
cd ../web
npm install
npm run build

# 2. Запустить gateway (будет serve статику из ../web/dist)
cd ../admin-gateway
go run cmd/admin-gateway/main.go

# Открыть http://localhost:8080
```

### Docker

```bash
docker build -t admin-gateway .
docker run -p 8080:8080 \
  -e GRPC_VENUE_ADDR=venue-svc:50051 \
  -e GRPC_BOOKING_ADDR=booking-svc:50052 \
  -e REDIS_ADDR=redis-master:6379 \
  -v $(pwd)/../web/dist:/root/web/dist:ro \
  admin-gateway
```

### Docker Compose

```bash
cd ../infra
docker compose --profile infra-min --profile apps up -d admin-gateway
```

---

## 🌐 CORS

CORS включен для всех origin'ов (для разработки):

```go
e.Use(middleware.CORS())
```

**Для production** нужно ограничить allowed origins:

```go
e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowOrigins: []string{"https://yourdomain.com"},
    AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
}))
```

---

## 🐛 Troubleshooting

### Проблема: 401 Unauthorized

```bash
# Проверить формат токена
curl -v http://localhost:18080/api/venues \
  -H "Authorization: Bearer token-admin"

# Проверить что пользователь существует в Redis
redis-cli -h localhost -p 7379 -a redis_pass get user:admin
```

### Проблема: 429 Rate Limit Exceeded

```bash
# Посмотреть счетчик в Redis
redis-cli -h localhost -p 7379 -a redis_pass get rl:admin-1

# Сбросить счетчик
redis-cli -h localhost -p 7379 -a redis_pass del rl:admin-1
```

### Проблема: gRPC connection failed

```bash
# Проверить доступность venue-svc
grpcurl -plaintext localhost:50151 list

# Проверить доступность booking-svc
grpcurl -plaintext localhost:50152 list

# Проверить логи
docker compose logs venue-svc
docker compose logs booking-svc
```

### Проблема: Static files not found

```bash
# Проверить что фронтенд собран
ls ../web/dist/

# Проверить mount в docker-compose.yml
# Убедиться что volume правильно настроен
```

---

## 🔐 Security Considerations

### ⚠️ Для Production

Текущая реализация **НЕ БЕЗОПАСНА для production**. Необходимо:

1. **JWT вместо простых токенов**
```go
import "github.com/golang-jwt/jwt/v5"

func generateJWT(username string) (string, error) {
    claims := jwt.MapClaims{
        "username": username,
        "exp":      time.Now().Add(time.Hour * 24).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(cfg.JWTSecret))
}
```

2. **Хэширование паролей**
```go
import "golang.org/x/crypto/bcrypt"

func hashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

func checkPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

3. **PostgreSQL для пользователей** (не Redis)

4. **HTTPS** для production

5. **Refresh tokens** с rotation

6. **CSRF protection**

7. **Input validation** и sanitization

---

## 📚 Дополнительные ресурсы

- [Echo Framework](https://echo.labstack.com/) - документация
- [gRPC-Go](https://grpc.io/docs/languages/go/) - gRPC client
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

---

## 🔗 API Documentation

См. [API.md](./API.md) для полной документации API endpoints (TODO: создать).

### Примеры запросов

**Регистрация:**
```bash
curl -X POST http://localhost:18080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123",
    "email": "admin@example.com"
  }'
```

**Создание бронирования:**
```bash
curl -X POST http://localhost:18080/api/bookings \
  -H "Authorization: Bearer token-admin" \
  -H "Content-Type: application/json" \
  -d '{
    "venue_id": "venue-1",
    "table": {"table_id": "table-1"},
    "slot": {"date": "2024-12-25", "start_time": "19:00", "duration_minutes": 120},
    "party_size": 4,
    "customer_name": "Иван Иванов",
    "customer_phone": "+79991234567"
  }'
```
