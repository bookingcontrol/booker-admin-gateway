package booking

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	bookingpb "github.com/bookingcontrol/booker-contracts-go/booking"
)

// MockBookingRepository is a mock implementation of booking repository
type MockBookingRepository struct {
	mock.Mock
}

func (m *MockBookingRepository) ListBookings(ctx context.Context, req *bookingpb.ListBookingsRequest) (*bookingpb.ListBookingsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookingpb.ListBookingsResponse), args.Error(1)
}

func (m *MockBookingRepository) GetBooking(ctx context.Context, id string) (*bookingpb.Booking, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookingpb.Booking), args.Error(1)
}

func (m *MockBookingRepository) CreateBooking(ctx context.Context, req *bookingpb.CreateBookingRequest) (*bookingpb.Booking, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookingpb.Booking), args.Error(1)
}

func (m *MockBookingRepository) ConfirmBooking(ctx context.Context, id, adminID string) (*bookingpb.Booking, error) {
	args := m.Called(ctx, id, adminID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookingpb.Booking), args.Error(1)
}

func (m *MockBookingRepository) CancelBooking(ctx context.Context, id, adminID, reason string) (*bookingpb.Booking, error) {
	args := m.Called(ctx, id, adminID, reason)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookingpb.Booking), args.Error(1)
}

func (m *MockBookingRepository) MarkSeated(ctx context.Context, id, adminID string) (*bookingpb.Booking, error) {
	args := m.Called(ctx, id, adminID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookingpb.Booking), args.Error(1)
}

func (m *MockBookingRepository) MarkFinished(ctx context.Context, id, adminID string) (*bookingpb.Booking, error) {
	args := m.Called(ctx, id, adminID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookingpb.Booking), args.Error(1)
}

func (m *MockBookingRepository) MarkNoShow(ctx context.Context, id, adminID string) (*bookingpb.Booking, error) {
	args := m.Called(ctx, id, adminID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookingpb.Booking), args.Error(1)
}

func TestService_ListBookings(t *testing.T) {
	t.Run("successful list", func(t *testing.T) {
		mockRepo := new(MockBookingRepository)
		service := NewService(mockRepo)

		req := &bookingpb.ListBookingsRequest{
			VenueId: "venue-1",
			Limit:   50,
			Offset:  0,
		}
		expected := &bookingpb.ListBookingsResponse{
			Bookings: []*bookingpb.Booking{
				{Id: "booking-1", VenueId: "venue-1"},
			},
			Total: 1,
		}
		mockRepo.On("ListBookings", mock.Anything, req).Return(expected, nil)

		result, err := service.ListBookings(context.Background(), req)

		require.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockBookingRepository)
		service := NewService(mockRepo)

		req := &bookingpb.ListBookingsRequest{VenueId: "venue-1"}
		mockRepo.On("ListBookings", mock.Anything, req).Return(nil, errors.New("db error"))

		_, err := service.ListBookings(context.Background(), req)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestService_GetBooking(t *testing.T) {
	t.Run("successful get", func(t *testing.T) {
		mockRepo := new(MockBookingRepository)
		service := NewService(mockRepo)

		expected := &bookingpb.Booking{Id: "booking-1", VenueId: "venue-1", Status: "confirmed"}
		mockRepo.On("GetBooking", mock.Anything, "booking-1").Return(expected, nil)

		result, err := service.GetBooking(context.Background(), "booking-1")

		require.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("booking not found", func(t *testing.T) {
		mockRepo := new(MockBookingRepository)
		service := NewService(mockRepo)

		mockRepo.On("GetBooking", mock.Anything, "nonexistent").Return(nil, errors.New("not found"))

		_, err := service.GetBooking(context.Background(), "nonexistent")

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestService_CreateBooking(t *testing.T) {
	t.Run("successful create", func(t *testing.T) {
		mockRepo := new(MockBookingRepository)
		service := NewService(mockRepo)

		req := &bookingpb.CreateBookingRequest{
			VenueId:      "venue-1",
			PartySize:    4,
			CustomerName: "John Doe",
		}
		expected := &bookingpb.Booking{
			Id:            "booking-new",
			VenueId:       "venue-1",
			PartySize:     4,
			CustomerName:  "John Doe",
			Status:        "held",
		}
		mockRepo.On("CreateBooking", mock.Anything, req).Return(expected, nil)

		result, err := service.CreateBooking(context.Background(), req)

		require.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestService_ConfirmBooking(t *testing.T) {
	t.Run("successful confirm", func(t *testing.T) {
		mockRepo := new(MockBookingRepository)
		service := NewService(mockRepo)

		expected := &bookingpb.Booking{Id: "booking-1", Status: "confirmed"}
		mockRepo.On("ConfirmBooking", mock.Anything, "booking-1", "admin-1").Return(expected, nil)

		result, err := service.ConfirmBooking(context.Background(), "booking-1", "admin-1")

		require.NoError(t, err)
		assert.Equal(t, "confirmed", result.Status)
		mockRepo.AssertExpectations(t)
	})
}

func TestService_CancelBooking(t *testing.T) {
	t.Run("successful cancel", func(t *testing.T) {
		mockRepo := new(MockBookingRepository)
		service := NewService(mockRepo)

		expected := &bookingpb.Booking{Id: "booking-1", Status: "cancelled"}
		mockRepo.On("CancelBooking", mock.Anything, "booking-1", "admin-1", "No show").Return(expected, nil)

		result, err := service.CancelBooking(context.Background(), "booking-1", "admin-1", "No show")

		require.NoError(t, err)
		assert.Equal(t, "cancelled", result.Status)
		mockRepo.AssertExpectations(t)
	})
}

func TestService_MarkSeated(t *testing.T) {
	t.Run("successful mark seated", func(t *testing.T) {
		mockRepo := new(MockBookingRepository)
		service := NewService(mockRepo)

		expected := &bookingpb.Booking{Id: "booking-1", Status: "seated"}
		mockRepo.On("MarkSeated", mock.Anything, "booking-1", "admin-1").Return(expected, nil)

		result, err := service.MarkSeated(context.Background(), "booking-1", "admin-1")

		require.NoError(t, err)
		assert.Equal(t, "seated", result.Status)
		mockRepo.AssertExpectations(t)
	})
}

func TestService_MarkFinished(t *testing.T) {
	t.Run("successful mark finished", func(t *testing.T) {
		mockRepo := new(MockBookingRepository)
		service := NewService(mockRepo)

		expected := &bookingpb.Booking{Id: "booking-1", Status: "finished"}
		mockRepo.On("MarkFinished", mock.Anything, "booking-1", "admin-1").Return(expected, nil)

		result, err := service.MarkFinished(context.Background(), "booking-1", "admin-1")

		require.NoError(t, err)
		assert.Equal(t, "finished", result.Status)
		mockRepo.AssertExpectations(t)
	})
}

func TestService_MarkNoShow(t *testing.T) {
	t.Run("successful mark no show", func(t *testing.T) {
		mockRepo := new(MockBookingRepository)
		service := NewService(mockRepo)

		expected := &bookingpb.Booking{Id: "booking-1", Status: "no_show"}
		mockRepo.On("MarkNoShow", mock.Anything, "booking-1", "admin-1").Return(expected, nil)

		result, err := service.MarkNoShow(context.Background(), "booking-1", "admin-1")

		require.NoError(t, err)
		assert.Equal(t, "no_show", result.Status)
		mockRepo.AssertExpectations(t)
	})
}

