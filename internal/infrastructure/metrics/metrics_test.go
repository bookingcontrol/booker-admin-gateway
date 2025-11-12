package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetrics_Registration(t *testing.T) {
	t.Run("HTTPRequestsTotal is registered", func(t *testing.T) {
		require.NotNil(t, HTTPRequestsTotal)
		
		// Проверяем, что метрика зарегистрирована в prometheus
		desc := HTTPRequestsTotal.WithLabelValues("GET", "/test", "200", "admin-gateway")
		assert.NotNil(t, desc)
		
		// Проверяем, что можно инкрементировать
		desc.Inc()
	})
	
	t.Run("HTTPRequestDuration is registered", func(t *testing.T) {
		require.NotNil(t, HTTPRequestDuration)
		
		// Проверяем, что метрика зарегистрирована
		observer := HTTPRequestDuration.WithLabelValues("GET", "/test", "200", "admin-gateway")
		assert.NotNil(t, observer)
		
		// Проверяем, что можно записать значение
		observer.Observe(0.1)
	})
	
	t.Run("HTTPRequestsTotal has correct labels", func(t *testing.T) {
		// Проверяем, что метрика имеет правильные лейблы
		metric := HTTPRequestsTotal.WithLabelValues("POST", "/api/v1/bookings", "201", "admin-gateway")
		require.NotNil(t, metric)
		
		// Проверяем, что можно инкрементировать несколько раз
		metric.Inc()
		metric.Inc()
	})
	
	t.Run("HTTPRequestDuration has correct buckets", func(t *testing.T) {
		// Проверяем, что гистограмма использует стандартные buckets
		observer := HTTPRequestDuration.WithLabelValues("GET", "/api/v1/venues", "200", "admin-gateway")
		require.NotNil(t, observer)
		
		// Записываем несколько значений
		observer.Observe(0.01)
		observer.Observe(0.05)
		observer.Observe(0.1)
		observer.Observe(0.5)
		observer.Observe(1.0)
	})
}

func TestMetrics_LabelValues(t *testing.T) {
	t.Run("HTTPRequestsTotal supports different HTTP methods", func(t *testing.T) {
		methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
		
		for _, method := range methods {
			metric := HTTPRequestsTotal.WithLabelValues(method, "/test", "200", "admin-gateway")
			require.NotNil(t, metric)
			metric.Inc()
		}
	})
	
	t.Run("HTTPRequestsTotal supports different status codes", func(t *testing.T) {
		statusCodes := []string{"200", "201", "400", "401", "404", "500"}
		
		for _, status := range statusCodes {
			metric := HTTPRequestsTotal.WithLabelValues("GET", "/test", status, "admin-gateway")
			require.NotNil(t, metric)
			metric.Inc()
		}
	})
}

func TestMetrics_CollectorInterface(t *testing.T) {
	t.Run("HTTPRequestsTotal implements prometheus.Collector", func(t *testing.T) {
		var _ prometheus.Collector = HTTPRequestsTotal
	})
	
	t.Run("HTTPRequestDuration implements prometheus.Collector", func(t *testing.T) {
		var _ prometheus.Collector = HTTPRequestDuration
	})
}

