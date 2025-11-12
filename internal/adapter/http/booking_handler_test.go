package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	bookingpb "github.com/bookingcontrol/booker-contracts-go/booking"
	uc "github.com/bookingcontrol/booker-admin-gateway/internal/usecase/booking"
)

// MockBookingRepository is a mock for booking repository
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

func TestBookingHandler_CreateBooking(t *testing.T) {
	e := echo.New()

	t.Run("successful create", func(t *testing.T) {
		mockRepo := new(MockBookingRepository)
		svc := uc.NewService(mockRepo)
		handler := NewBookingHandler(svc)

		reqBody := map[string]interface{}{
			"venue_id": "venue-1",
			"table": map[string]interface{}{
				"venue_id": "venue-1",
				"room_id":  "room-1",
				"table_id": "table-1",
			},
			"slot": map[string]interface{}{
				"date":              "2025-11-12",
				"start_time":        "18:00",
				"duration_minutes":  120,
			},
			"party_size":     2,
			"customer_name":  "John Doe",
			"customer_phone": "+79111111111",
			"comment":        "",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/bookings", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("admin_id", "admin-1")

		expected := &bookingpb.Booking{
			Id:            "booking-new",
			VenueId:       "venue-1",
			PartySize:     2,
			CustomerName:  "John Doe",
			CustomerPhone: "+79111111111",
			Status:        "held",
		}
		mockRepo.On("CreateBooking", mock.Anything, mock.MatchedBy(func(r *bookingpb.CreateBookingRequest) bool {
			return r.VenueId == "venue-1" && r.PartySize == 2 && r.AdminId == "admin-1"
		})).Return(expected, nil)

		err := handler.CreateBooking(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid request body", func(t *testing.T) {
		mockRepo := new(MockBookingRepository)
		svc := uc.NewService(mockRepo)
		handler := NewBookingHandler(svc)

		req := httptest.NewRequest(http.MethodPost, "/bookings", bytes.NewReader([]byte("invalid json")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.CreateBooking(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockRepo.AssertNotCalled(t, "CreateBooking")
	})
}

func TestBookingHandler_GetBooking(t *testing.T) {
	e := echo.New()

	t.Run("successful get", func(t *testing.T) {
		mockRepo := new(MockBookingRepository)
		svc := uc.NewService(mockRepo)
		handler := NewBookingHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/bookings/booking-1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/bookings/:id")
		c.SetParamNames("id")
		c.SetParamValues("booking-1")

		expected := &bookingpb.Booking{Id: "booking-1", VenueId: "venue-1", Status: "confirmed"}
		mockRepo.On("GetBooking", mock.Anything, "booking-1").Return(expected, nil)

		err := handler.GetBooking(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("booking not found", func(t *testing.T) {
		mockRepo := new(MockBookingRepository)
		svc := uc.NewService(mockRepo)
		handler := NewBookingHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/bookings/nonexistent", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/bookings/:id")
		c.SetParamNames("id")
		c.SetParamValues("nonexistent")

		mockRepo.On("GetBooking", mock.Anything, "nonexistent").Return(nil, errors.New("not found"))

		err := handler.GetBooking(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestBookingHandler_ConfirmBooking(t *testing.T) {
	e := echo.New()

	t.Run("successful confirm", func(t *testing.T) {
		mockRepo := new(MockBookingRepository)
		svc := uc.NewService(mockRepo)
		handler := NewBookingHandler(svc)

		req := httptest.NewRequest(http.MethodPost, "/bookings/booking-1/confirm", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/bookings/:id/confirm")
		c.SetParamNames("id")
		c.SetParamValues("booking-1")
		c.Set("admin_id", "admin-1")

		expected := &bookingpb.Booking{Id: "booking-1", Status: "confirmed"}
		mockRepo.On("ConfirmBooking", mock.Anything, "booking-1", "admin-1").Return(expected, nil)

		err := handler.ConfirmBooking(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestBookingHandler_CancelBooking(t *testing.T) {
	e := echo.New()

	t.Run("successful cancel", func(t *testing.T) {
		mockRepo := new(MockBookingRepository)
		svc := uc.NewService(mockRepo)
		handler := NewBookingHandler(svc)

		reqBody := map[string]interface{}{"reason": "Customer cancelled"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/bookings/booking-1/cancel", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/bookings/:id/cancel")
		c.SetParamNames("id")
		c.SetParamValues("booking-1")
		c.Set("admin_id", "admin-1")

		expected := &bookingpb.Booking{Id: "booking-1", Status: "cancelled"}
		mockRepo.On("CancelBooking", mock.Anything, "booking-1", "admin-1", "Customer cancelled").Return(expected, nil)

		err := handler.CancelBooking(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestBookingHandler_ListBookings(t *testing.T) {
	e := echo.New()

	t.Run("successful list", func(t *testing.T) {
		mockRepo := new(MockBookingRepository)
		svc := uc.NewService(mockRepo)
		handler := NewBookingHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/bookings?venue_id=venue-1&limit=50&offset=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		expected := &bookingpb.ListBookingsResponse{
			Bookings: []*bookingpb.Booking{{Id: "booking-1", VenueId: "venue-1"}},
			Total:    1,
		}
		mockRepo.On("ListBookings", mock.Anything, mock.MatchedBy(func(r *bookingpb.ListBookingsRequest) bool {
			return r.VenueId == "venue-1" && r.Limit == 50
		})).Return(expected, nil)

		err := handler.ListBookings(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("default limit when not provided", func(t *testing.T) {
		mockRepo := new(MockBookingRepository)
		svc := uc.NewService(mockRepo)
		handler := NewBookingHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/bookings", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		expected := &bookingpb.ListBookingsResponse{Bookings: []*bookingpb.Booking{}, Total: 0}
		mockRepo.On("ListBookings", mock.Anything, mock.MatchedBy(func(r *bookingpb.ListBookingsRequest) bool {
			return r.Limit == 50
		})).Return(expected, nil)

		err := handler.ListBookings(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestBookingHandler_MarkSeated(t *testing.T) {
	e := echo.New()

	t.Run("successful mark seated", func(t *testing.T) {
		mockRepo := new(MockBookingRepository)
		svc := uc.NewService(mockRepo)
		handler := NewBookingHandler(svc)

		req := httptest.NewRequest(http.MethodPost, "/bookings/booking-1/seat", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/bookings/:id/seat")
		c.SetParamNames("id")
		c.SetParamValues("booking-1")
		c.Set("admin_id", "admin-1")

		expected := &bookingpb.Booking{Id: "booking-1", Status: "seated"}
		mockRepo.On("MarkSeated", mock.Anything, "booking-1", "admin-1").Return(expected, nil)

		err := handler.MarkSeated(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestBookingHandler_MarkFinished(t *testing.T) {
	e := echo.New()

	t.Run("successful mark finished", func(t *testing.T) {
		mockRepo := new(MockBookingRepository)
		svc := uc.NewService(mockRepo)
		handler := NewBookingHandler(svc)

		req := httptest.NewRequest(http.MethodPost, "/bookings/booking-1/finish", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/bookings/:id/finish")
		c.SetParamNames("id")
		c.SetParamValues("booking-1")
		c.Set("admin_id", "admin-1")

		expected := &bookingpb.Booking{Id: "booking-1", Status: "finished"}
		mockRepo.On("MarkFinished", mock.Anything, "booking-1", "admin-1").Return(expected, nil)

		err := handler.MarkFinished(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestBookingHandler_MarkNoShow(t *testing.T) {
	e := echo.New()

	t.Run("successful mark no show", func(t *testing.T) {
		mockRepo := new(MockBookingRepository)
		svc := uc.NewService(mockRepo)
		handler := NewBookingHandler(svc)

		req := httptest.NewRequest(http.MethodPost, "/bookings/booking-1/no-show", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/bookings/:id/no-show")
		c.SetParamNames("id")
		c.SetParamValues("booking-1")
		c.Set("admin_id", "admin-1")

		expected := &bookingpb.Booking{Id: "booking-1", Status: "no_show"}
		mockRepo.On("MarkNoShow", mock.Anything, "booking-1", "admin-1").Return(expected, nil)

		err := handler.MarkNoShow(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

