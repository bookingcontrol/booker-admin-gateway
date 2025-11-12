package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	bookingpb "github.com/bookingcontrol/booker-contracts-go/booking"
	venuepb "github.com/bookingcontrol/booker-contracts-go/venue"
	ucauth "github.com/bookingcontrol/booker-admin-gateway/internal/usecase/auth"
	ucbooking "github.com/bookingcontrol/booker-admin-gateway/internal/usecase/booking"
	ucvenue "github.com/bookingcontrol/booker-admin-gateway/internal/usecase/venue"
)

// Интеграционные тесты: проверяем полный flow handler -> use case -> repository (repository мокируется)

// MockAuthRepo для интеграционных тестов
type MockAuthRepoIntegration struct {
	mock.Mock
}

func (m *MockAuthRepoIntegration) UserExists(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthRepoIntegration) GetUserPassword(ctx context.Context, username string) (string, error) {
	args := m.Called(ctx, username)
	return args.String(0), args.Error(1)
}

func (m *MockAuthRepoIntegration) CreateUser(ctx context.Context, username string, userData map[string]interface{}) error {
	args := m.Called(ctx, username, userData)
	return args.Error(0)
}

// MockVenueRepo для интеграционных тестов
type MockVenueRepoIntegration struct {
	mock.Mock
}

func (m *MockVenueRepoIntegration) ListVenues(ctx context.Context, limit, offset int32) (*venuepb.ListVenuesResponse, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*venuepb.ListVenuesResponse), args.Error(1)
}

func (m *MockVenueRepoIntegration) GetVenue(ctx context.Context, id string) (*venuepb.Venue, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*venuepb.Venue), args.Error(1)
}

func (m *MockVenueRepoIntegration) CreateVenue(ctx context.Context, req *venuepb.CreateVenueRequest) (*venuepb.Venue, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*venuepb.Venue), args.Error(1)
}

func (m *MockVenueRepoIntegration) UpdateVenue(ctx context.Context, req *venuepb.UpdateVenueRequest) (*venuepb.Venue, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*venuepb.Venue), args.Error(1)
}

func (m *MockVenueRepoIntegration) DeleteVenue(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockVenueRepoIntegration) ListRooms(ctx context.Context, venueID string, limit, offset int32) (*venuepb.ListRoomsResponse, error) {
	args := m.Called(ctx, venueID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*venuepb.ListRoomsResponse), args.Error(1)
}

func (m *MockVenueRepoIntegration) GetRoom(ctx context.Context, id string) (*venuepb.Room, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*venuepb.Room), args.Error(1)
}

func (m *MockVenueRepoIntegration) CreateRoom(ctx context.Context, req *venuepb.CreateRoomRequest) (*venuepb.Room, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*venuepb.Room), args.Error(1)
}

func (m *MockVenueRepoIntegration) UpdateRoom(ctx context.Context, req *venuepb.UpdateRoomRequest) (*venuepb.Room, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*venuepb.Room), args.Error(1)
}

func (m *MockVenueRepoIntegration) DeleteRoom(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockVenueRepoIntegration) ListTables(ctx context.Context, roomID string, limit, offset int32) (*venuepb.ListTablesResponse, error) {
	args := m.Called(ctx, roomID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*venuepb.ListTablesResponse), args.Error(1)
}

func (m *MockVenueRepoIntegration) GetTable(ctx context.Context, id string) (*venuepb.Table, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*venuepb.Table), args.Error(1)
}

func (m *MockVenueRepoIntegration) CreateTable(ctx context.Context, req *venuepb.CreateTableRequest) (*venuepb.Table, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*venuepb.Table), args.Error(1)
}

func (m *MockVenueRepoIntegration) UpdateTable(ctx context.Context, req *venuepb.UpdateTableRequest) (*venuepb.Table, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*venuepb.Table), args.Error(1)
}

func (m *MockVenueRepoIntegration) DeleteTable(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockVenueRepoIntegration) GetOpeningHours(ctx context.Context, venueID string) (*venuepb.OpeningHours, error) {
	args := m.Called(ctx, venueID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*venuepb.OpeningHours), args.Error(1)
}

func (m *MockVenueRepoIntegration) SetOpeningHours(ctx context.Context, req *venuepb.SetOpeningHoursRequest) (*venuepb.SetOpeningHoursResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*venuepb.SetOpeningHoursResponse), args.Error(1)
}

func (m *MockVenueRepoIntegration) SetSpecialHours(ctx context.Context, req *venuepb.SetSpecialHoursRequest) (*venuepb.SetSpecialHoursResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*venuepb.SetSpecialHoursResponse), args.Error(1)
}

func (m *MockVenueRepoIntegration) CheckAvailability(ctx context.Context, req *venuepb.CheckAvailabilityRequest) (*venuepb.CheckAvailabilityResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*venuepb.CheckAvailabilityResponse), args.Error(1)
}

// MockBookingRepo для интеграционных тестов
type MockBookingRepoIntegration struct {
	mock.Mock
}

func (m *MockBookingRepoIntegration) ListBookings(ctx context.Context, req *bookingpb.ListBookingsRequest) (*bookingpb.ListBookingsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookingpb.ListBookingsResponse), args.Error(1)
}

func (m *MockBookingRepoIntegration) GetBooking(ctx context.Context, id string) (*bookingpb.Booking, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookingpb.Booking), args.Error(1)
}

func (m *MockBookingRepoIntegration) CreateBooking(ctx context.Context, req *bookingpb.CreateBookingRequest) (*bookingpb.Booking, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookingpb.Booking), args.Error(1)
}

func (m *MockBookingRepoIntegration) ConfirmBooking(ctx context.Context, id, adminID string) (*bookingpb.Booking, error) {
	args := m.Called(ctx, id, adminID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookingpb.Booking), args.Error(1)
}

func (m *MockBookingRepoIntegration) CancelBooking(ctx context.Context, id, adminID, reason string) (*bookingpb.Booking, error) {
	args := m.Called(ctx, id, adminID, reason)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookingpb.Booking), args.Error(1)
}

func (m *MockBookingRepoIntegration) MarkSeated(ctx context.Context, id, adminID string) (*bookingpb.Booking, error) {
	args := m.Called(ctx, id, adminID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookingpb.Booking), args.Error(1)
}

func (m *MockBookingRepoIntegration) MarkFinished(ctx context.Context, id, adminID string) (*bookingpb.Booking, error) {
	args := m.Called(ctx, id, adminID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookingpb.Booking), args.Error(1)
}

func (m *MockBookingRepoIntegration) MarkNoShow(ctx context.Context, id, adminID string) (*bookingpb.Booking, error) {
	args := m.Called(ctx, id, adminID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookingpb.Booking), args.Error(1)
}

// Интеграционный тест: полный flow от HTTP запроса до repository
func TestIntegration_AuthFlow(t *testing.T) {
	e := echo.New()
	
	// Создаем реальную цепочку: handler -> use case -> repository (мок)
	mockAuthRepo := new(MockAuthRepoIntegration)
	authSvc := ucauth.NewService(mockAuthRepo)
	authHandler := NewAuthHandler(authSvc)
	
	t.Run("full registration flow", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"username": "newuser",
			"password": "password123",
			"email":    "new@example.com",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		
		// Мокаем repository (внешняя зависимость)
		mockAuthRepo.On("UserExists", mock.Anything, "newuser").Return(false, nil)
		mockAuthRepo.On("CreateUser", mock.Anything, "newuser", mock.AnythingOfType("map[string]interface {}")).Return(nil)
		
		// Вызываем handler (который вызывает use case, который вызывает repository)
		err := authHandler.Register(c)
		
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		
		// Проверяем, что все слои были вызваны
		mockAuthRepo.AssertExpectations(t)
		
		// Проверяем ответ
		var response ucauth.RegisterView
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, "newuser", response.Username)
	})
	
	t.Run("full login flow", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"username": "testuser",
			"password": "password123",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		
		// Мокаем repository
		mockAuthRepo.On("UserExists", mock.Anything, "testuser").Return(true, nil)
		mockAuthRepo.On("GetUserPassword", mock.Anything, "testuser").Return("password123", nil)
		
		err := authHandler.Login(c)
		
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockAuthRepo.AssertExpectations(t)
		
		var response ucauth.LoginView
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NotEmpty(t, response.AccessToken)
	})
}

func TestIntegration_BookingFlow(t *testing.T) {
	e := echo.New()
	
	// Создаем реальную цепочку: handler -> use case -> repository (мок)
	mockBookingRepo := new(MockBookingRepoIntegration)
	bookingSvc := ucbooking.NewService(mockBookingRepo)
	bookingHandler := NewBookingHandler(bookingSvc)
	
	t.Run("full create booking flow", func(t *testing.T) {
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
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/bookings", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("admin_id", "admin-1")
		
		// Мокаем repository (gRPC вызов)
		expectedBooking := &bookingpb.Booking{
			Id:            "booking-new",
			VenueId:       "venue-1",
			PartySize:     2,
			CustomerName:  "John Doe",
			CustomerPhone: "+79111111111",
			Status:        "held",
		}
		mockBookingRepo.On("CreateBooking", mock.Anything, mock.MatchedBy(func(r *bookingpb.CreateBookingRequest) bool {
			return r.VenueId == "venue-1" && r.PartySize == 2 && r.AdminId == "admin-1"
		})).Return(expectedBooking, nil)
		
		// Вызываем handler -> use case -> repository
		err := bookingHandler.CreateBooking(c)
		
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		mockBookingRepo.AssertExpectations(t)
		
		// Проверяем, что данные прошли через все слои
		var response bookingpb.Booking
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, "booking-new", response.Id)
		assert.Equal(t, "venue-1", response.VenueId)
	})
	
	t.Run("full confirm booking flow", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/bookings/booking-1/confirm", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/bookings/:id/confirm")
		c.SetParamNames("id")
		c.SetParamValues("booking-1")
		c.Set("admin_id", "admin-1")
		
		expectedBooking := &bookingpb.Booking{Id: "booking-1", Status: "confirmed"}
		mockBookingRepo.On("ConfirmBooking", mock.Anything, "booking-1", "admin-1").Return(expectedBooking, nil)
		
		err := bookingHandler.ConfirmBooking(c)
		
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockBookingRepo.AssertExpectations(t)
	})
}

func TestIntegration_VenueFlow(t *testing.T) {
	e := echo.New()
	
	// Создаем реальную цепочку: handler -> use case -> repository (мок)
	mockVenueRepo := new(MockVenueRepoIntegration)
	venueSvc := ucvenue.NewService(mockVenueRepo)
	venueHandler := NewVenueHandler(venueSvc)
	
	t.Run("full create venue flow", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name":     "New Venue",
			"timezone": "UTC",
			"address":  "123 Main St",
			"phone":    "+1234567890",
			"email":    "venue@example.com",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/venues", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		
		// Мокаем repository (gRPC вызов)
		expectedVenue := &venuepb.Venue{
			Id:       "venue-new",
			Name:     "New Venue",
			Timezone: "UTC",
			Address:  "123 Main St",
			Phone:    "+1234567890",
			Email:    "venue@example.com",
		}
		mockVenueRepo.On("CreateVenue", mock.Anything, mock.MatchedBy(func(r *venuepb.CreateVenueRequest) bool {
			return r.Name == "New Venue" && r.Timezone == "UTC"
		})).Return(expectedVenue, nil)
		
		// Вызываем handler -> use case -> repository
		err := venueHandler.CreateVenue(c)
		
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		mockVenueRepo.AssertExpectations(t)
		
		var response venuepb.Venue
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, "venue-new", response.Id)
		assert.Equal(t, "New Venue", response.Name)
	})
	
	t.Run("full list venues flow", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/venues?limit=10&offset=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		
		expected := &venuepb.ListVenuesResponse{
			Venues: []*venuepb.Venue{
				{Id: "venue-1", Name: "Venue 1"},
				{Id: "venue-2", Name: "Venue 2"},
			},
			Total: 2,
		}
		mockVenueRepo.On("ListVenues", mock.Anything, int32(10), int32(0)).Return(expected, nil)
		
		err := venueHandler.ListVenues(c)
		
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockVenueRepo.AssertExpectations(t)
	})
}

