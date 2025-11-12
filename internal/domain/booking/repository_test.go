package booking

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	bookingpb "github.com/bookingcontrol/booker-contracts-go/booking"
)

// Тестируем контракт интерфейса Repository
// Проверяем, что интерфейс правильно определен

// MockRepository - пример реализации для тестирования контракта
type MockRepository struct {
	GetBookingFunc func(ctx context.Context, id string) (*bookingpb.Booking, error)
}

func (m *MockRepository) ListBookings(ctx context.Context, req *bookingpb.ListBookingsRequest) (*bookingpb.ListBookingsResponse, error) {
	return nil, nil
}

func (m *MockRepository) GetBooking(ctx context.Context, id string) (*bookingpb.Booking, error) {
	if m.GetBookingFunc != nil {
		return m.GetBookingFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockRepository) CreateBooking(ctx context.Context, req *bookingpb.CreateBookingRequest) (*bookingpb.Booking, error) {
	return nil, nil
}

func (m *MockRepository) ConfirmBooking(ctx context.Context, id, adminID string) (*bookingpb.Booking, error) {
	return nil, nil
}

func (m *MockRepository) CancelBooking(ctx context.Context, id, adminID, reason string) (*bookingpb.Booking, error) {
	return nil, nil
}

func (m *MockRepository) MarkSeated(ctx context.Context, id, adminID string) (*bookingpb.Booking, error) {
	return nil, nil
}

func (m *MockRepository) MarkFinished(ctx context.Context, id, adminID string) (*bookingpb.Booking, error) {
	return nil, nil
}

func (m *MockRepository) MarkNoShow(ctx context.Context, id, adminID string) (*bookingpb.Booking, error) {
	return nil, nil
}

func TestRepositoryInterface(t *testing.T) {
	t.Run("MockRepository implements Repository interface", func(t *testing.T) {
		var _ Repository = (*MockRepository)(nil)
	})
	
	t.Run("Repository methods have correct signatures", func(t *testing.T) {
		repo := &MockRepository{
			GetBookingFunc: func(ctx context.Context, id string) (*bookingpb.Booking, error) {
				return &bookingpb.Booking{
					Id:     id,
					Status: "confirmed",
				}, nil
			},
		}
		
		ctx := context.Background()
		booking, err := repo.GetBooking(ctx, "booking-1")
		
		assert.NoError(t, err)
		assert.NotNil(t, booking)
		assert.Equal(t, "booking-1", booking.Id)
		assert.Equal(t, "confirmed", booking.Status)
	})
	
	t.Run("ConfirmBooking signature", func(t *testing.T) {
		repo := &MockRepository{}
		ctx := context.Background()
		
		// Проверяем, что метод принимает правильные параметры
		_, err := repo.ConfirmBooking(ctx, "booking-1", "admin-1")
		assert.NoError(t, err)
	})
	
	t.Run("CancelBooking signature", func(t *testing.T) {
		repo := &MockRepository{}
		ctx := context.Background()
		
		// Проверяем, что метод принимает правильные параметры
		_, err := repo.CancelBooking(ctx, "booking-1", "admin-1", "reason")
		assert.NoError(t, err)
	})
}

