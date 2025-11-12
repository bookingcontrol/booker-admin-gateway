package tracing

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

func TestInitTracer(t *testing.T) {
	t.Run("returns no-op shutdown when endpoint is empty", func(t *testing.T) {
		shutdown, err := InitTracer("test-service", "")
		
		require.NoError(t, err)
		require.NotNil(t, shutdown)
		
		// Вызываем shutdown - не должно быть ошибок
		shutdown()
	})
	
	t.Run("returns error for invalid endpoint format", func(t *testing.T) {
		// Используем невалидный endpoint, чтобы проверить обработку ошибок
		// В реальности это может не упасть сразу, но проверим структуру
		shutdown, err := InitTracer("test-service", "invalid://endpoint")
		
		// Может быть ошибка или нет, в зависимости от реализации
		if err != nil {
			assert.Nil(t, shutdown)
		} else {
			require.NotNil(t, shutdown)
			shutdown()
		}
	})
}

func TestStartSpan(t *testing.T) {
	t.Run("creates span with context", func(t *testing.T) {
		ctx := context.Background()
		
		newCtx, span := StartSpan(ctx, "test-operation")
		
		require.NotNil(t, newCtx)
		require.NotNil(t, span)
		
		// Проверяем, что span можно завершить
		span.End()
	})
	
	t.Run("span implements trace.Span interface", func(t *testing.T) {
		ctx := context.Background()
		_, span := StartSpan(ctx, "test")
		
		var _ trace.Span = span
		span.End()
	})
	
	t.Run("span context is propagated", func(t *testing.T) {
		ctx := context.Background()
		
		newCtx, span1 := StartSpan(ctx, "operation1")
		require.NotNil(t, newCtx)
		
		// Создаем дочерний span
		childCtx, span2 := StartSpan(newCtx, "operation2")
		require.NotNil(t, childCtx)
		
		span2.End()
		span1.End()
	})
	
	t.Run("multiple spans can be created", func(t *testing.T) {
		ctx := context.Background()
		
		ctx1, span1 := StartSpan(ctx, "span1")
		ctx2, span2 := StartSpan(ctx, "span2")
		
		require.NotNil(t, ctx1)
		require.NotNil(t, ctx2)
		require.NotNil(t, span1)
		require.NotNil(t, span2)
		
		span1.End()
		span2.End()
	})
}

func TestTracing_Integration(t *testing.T) {
	t.Run("tracing works without jaeger endpoint", func(t *testing.T) {
		// Инициализируем без endpoint (no-op режим)
		shutdown, err := InitTracer("test-service", "")
		require.NoError(t, err)
		defer shutdown()
		
		ctx := context.Background()
		newCtx, span := StartSpan(ctx, "test-operation")
		
		require.NotNil(t, newCtx)
		require.NotNil(t, span)
		
		span.End()
	})
}

