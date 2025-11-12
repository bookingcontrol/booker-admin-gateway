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
	venuepb "github.com/bookingcontrol/booker-contracts-go/venue"
	uc "github.com/bookingcontrol/booker-admin-gateway/internal/usecase/venue"
)

// MockVenueRepository is a mock for venue repository
type MockVenueRepository struct {
	mock.Mock
}

func (m *MockVenueRepository) ListVenues(ctx context.Context, limit, offset int32) (*venuepb.ListVenuesResponse, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*venuepb.ListVenuesResponse), args.Error(1)
}

func (m *MockVenueRepository) GetVenue(ctx context.Context, id string) (*venuepb.Venue, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*venuepb.Venue), args.Error(1)
}

func (m *MockVenueRepository) CreateVenue(ctx context.Context, req *venuepb.CreateVenueRequest) (*venuepb.Venue, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*venuepb.Venue), args.Error(1)
}

func (m *MockVenueRepository) UpdateVenue(ctx context.Context, req *venuepb.UpdateVenueRequest) (*venuepb.Venue, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*venuepb.Venue), args.Error(1)
}

func (m *MockVenueRepository) DeleteVenue(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockVenueRepository) ListRooms(ctx context.Context, venueID string, limit, offset int32) (*venuepb.ListRoomsResponse, error) {
	args := m.Called(ctx, venueID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*venuepb.ListRoomsResponse), args.Error(1)
}

func (m *MockVenueRepository) GetRoom(ctx context.Context, id string) (*venuepb.Room, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*venuepb.Room), args.Error(1)
}

func (m *MockVenueRepository) CreateRoom(ctx context.Context, req *venuepb.CreateRoomRequest) (*venuepb.Room, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*venuepb.Room), args.Error(1)
}

func (m *MockVenueRepository) UpdateRoom(ctx context.Context, req *venuepb.UpdateRoomRequest) (*venuepb.Room, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*venuepb.Room), args.Error(1)
}

func (m *MockVenueRepository) DeleteRoom(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockVenueRepository) ListTables(ctx context.Context, roomID string, limit, offset int32) (*venuepb.ListTablesResponse, error) {
	args := m.Called(ctx, roomID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*venuepb.ListTablesResponse), args.Error(1)
}

func (m *MockVenueRepository) GetTable(ctx context.Context, id string) (*venuepb.Table, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*venuepb.Table), args.Error(1)
}

func (m *MockVenueRepository) CreateTable(ctx context.Context, req *venuepb.CreateTableRequest) (*venuepb.Table, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*venuepb.Table), args.Error(1)
}

func (m *MockVenueRepository) UpdateTable(ctx context.Context, req *venuepb.UpdateTableRequest) (*venuepb.Table, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*venuepb.Table), args.Error(1)
}

func (m *MockVenueRepository) DeleteTable(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockVenueRepository) GetOpeningHours(ctx context.Context, venueID string) (*venuepb.OpeningHours, error) {
	args := m.Called(ctx, venueID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*venuepb.OpeningHours), args.Error(1)
}

func (m *MockVenueRepository) SetOpeningHours(ctx context.Context, req *venuepb.SetOpeningHoursRequest) (*venuepb.SetOpeningHoursResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*venuepb.SetOpeningHoursResponse), args.Error(1)
}

func (m *MockVenueRepository) SetSpecialHours(ctx context.Context, req *venuepb.SetSpecialHoursRequest) (*venuepb.SetSpecialHoursResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*venuepb.SetSpecialHoursResponse), args.Error(1)
}

func (m *MockVenueRepository) CheckAvailability(ctx context.Context, req *venuepb.CheckAvailabilityRequest) (*venuepb.CheckAvailabilityResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*venuepb.CheckAvailabilityResponse), args.Error(1)
}

func TestVenueHandler_ListVenues(t *testing.T) {
	e := echo.New()

	t.Run("successful list", func(t *testing.T) {
		mockRepo := new(MockVenueRepository)
		svc := uc.NewService(mockRepo)
		handler := NewVenueHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/venues?limit=50&offset=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		expected := &venuepb.ListVenuesResponse{
			Venues: []*venuepb.Venue{{Id: "venue-1", Name: "Test Venue"}},
			Total:  1,
		}
		mockRepo.On("ListVenues", mock.Anything, int32(50), int32(0)).Return(expected, nil)

		err := handler.ListVenues(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("default limit when not provided", func(t *testing.T) {
		mockRepo := new(MockVenueRepository)
		svc := uc.NewService(mockRepo)
		handler := NewVenueHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/venues", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		expected := &venuepb.ListVenuesResponse{Venues: []*venuepb.Venue{}, Total: 0}
		mockRepo.On("ListVenues", mock.Anything, int32(50), int32(0)).Return(expected, nil)

		err := handler.ListVenues(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockVenueRepository)
		svc := uc.NewService(mockRepo)
		handler := NewVenueHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/venues", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRepo.On("ListVenues", mock.Anything, int32(50), int32(0)).Return(nil, errors.New("db error"))

		err := handler.ListVenues(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestVenueHandler_GetVenue(t *testing.T) {
	e := echo.New()

	t.Run("successful get", func(t *testing.T) {
		mockRepo := new(MockVenueRepository)
		svc := uc.NewService(mockRepo)
		handler := NewVenueHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/venues/venue-1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/venues/:id")
		c.SetParamNames("id")
		c.SetParamValues("venue-1")

		expected := &venuepb.Venue{Id: "venue-1", Name: "Test Venue"}
		mockRepo.On("GetVenue", mock.Anything, "venue-1").Return(expected, nil)

		err := handler.GetVenue(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("venue not found", func(t *testing.T) {
		mockRepo := new(MockVenueRepository)
		svc := uc.NewService(mockRepo)
		handler := NewVenueHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/venues/nonexistent", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/venues/:id")
		c.SetParamNames("id")
		c.SetParamValues("nonexistent")

		mockRepo.On("GetVenue", mock.Anything, "nonexistent").Return(nil, errors.New("not found"))

		err := handler.GetVenue(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestVenueHandler_CreateVenue(t *testing.T) {
	e := echo.New()

	t.Run("successful create", func(t *testing.T) {
		mockRepo := new(MockVenueRepository)
		svc := uc.NewService(mockRepo)
		handler := NewVenueHandler(svc)

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

		expected := &venuepb.Venue{
			Id:       "venue-new",
			Name:     "New Venue",
			Timezone: "UTC",
			Address:  "123 Main St",
			Phone:    "+1234567890",
			Email:    "venue@example.com",
		}
		mockRepo.On("CreateVenue", mock.Anything, mock.MatchedBy(func(r *venuepb.CreateVenueRequest) bool {
			return r.Name == "New Venue" && r.Timezone == "UTC"
		})).Return(expected, nil)

		err := handler.CreateVenue(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestVenueHandler_DeleteVenue(t *testing.T) {
	e := echo.New()

	t.Run("successful delete", func(t *testing.T) {
		mockRepo := new(MockVenueRepository)
		svc := uc.NewService(mockRepo)
		handler := NewVenueHandler(svc)

		req := httptest.NewRequest(http.MethodDelete, "/venues/venue-1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/venues/:id")
		c.SetParamNames("id")
		c.SetParamValues("venue-1")

		mockRepo.On("DeleteVenue", mock.Anything, "venue-1").Return(nil)

		err := handler.DeleteVenue(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockVenueRepository)
		svc := uc.NewService(mockRepo)
		handler := NewVenueHandler(svc)

		req := httptest.NewRequest(http.MethodDelete, "/venues/venue-1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/venues/:id")
		c.SetParamNames("id")
		c.SetParamValues("venue-1")

		mockRepo.On("DeleteVenue", mock.Anything, "venue-1").Return(errors.New("db error"))

		err := handler.DeleteVenue(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestVenueHandler_UpdateVenue(t *testing.T) {
	e := echo.New()

	t.Run("successful update", func(t *testing.T) {
		mockRepo := new(MockVenueRepository)
		svc := uc.NewService(mockRepo)
		handler := NewVenueHandler(svc)

		reqBody := map[string]interface{}{
			"name":    "Updated Venue",
			"address": "456 New St",
			"phone":   "+9876543210",
			"email":   "updated@example.com",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPut, "/venues/venue-1", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/venues/:id")
		c.SetParamNames("id")
		c.SetParamValues("venue-1")

		expected := &venuepb.Venue{Id: "venue-1", Name: "Updated Venue"}
		mockRepo.On("UpdateVenue", mock.Anything, mock.MatchedBy(func(r *venuepb.UpdateVenueRequest) bool {
			return r.Id == "venue-1" && r.Name == "Updated Venue"
		})).Return(expected, nil)

		err := handler.UpdateVenue(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestVenueHandler_ListRooms(t *testing.T) {
	e := echo.New()

	t.Run("successful list", func(t *testing.T) {
		mockRepo := new(MockVenueRepository)
		svc := uc.NewService(mockRepo)
		handler := NewVenueHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/venues/venue-1/rooms?limit=50&offset=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/venues/:venueId/rooms")
		c.SetParamNames("venueId")
		c.SetParamValues("venue-1")

		expected := &venuepb.ListRoomsResponse{
			Rooms: []*venuepb.Room{{Id: "room-1", Name: "Main Room"}},
			Total: 1,
		}
		mockRepo.On("ListRooms", mock.Anything, "venue-1", int32(50), int32(0)).Return(expected, nil)

		err := handler.ListRooms(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestVenueHandler_CreateRoom(t *testing.T) {
	e := echo.New()

	t.Run("successful create", func(t *testing.T) {
		mockRepo := new(MockVenueRepository)
		svc := uc.NewService(mockRepo)
		handler := NewVenueHandler(svc)

		reqBody := map[string]interface{}{"name": "New Room"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/venues/venue-1/rooms", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/venues/:venueId/rooms")
		c.SetParamNames("venueId")
		c.SetParamValues("venue-1")

		expected := &venuepb.Room{Id: "room-new", Name: "New Room"}
		mockRepo.On("CreateRoom", mock.Anything, mock.MatchedBy(func(r *venuepb.CreateRoomRequest) bool {
			return r.VenueId == "venue-1" && r.Name == "New Room"
		})).Return(expected, nil)

		err := handler.CreateRoom(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestVenueHandler_ListTables(t *testing.T) {
	e := echo.New()

	t.Run("successful list", func(t *testing.T) {
		mockRepo := new(MockVenueRepository)
		svc := uc.NewService(mockRepo)
		handler := NewVenueHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/rooms/room-1/tables?limit=50&offset=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/rooms/:roomId/tables")
		c.SetParamNames("roomId")
		c.SetParamValues("room-1")

		expected := &venuepb.ListTablesResponse{
			Tables: []*venuepb.Table{{Id: "table-1", Name: "Table 1"}},
			Total:  1,
		}
		mockRepo.On("ListTables", mock.Anything, "room-1", int32(50), int32(0)).Return(expected, nil)

		err := handler.ListTables(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestVenueHandler_CreateTable(t *testing.T) {
	e := echo.New()

	t.Run("successful create", func(t *testing.T) {
		mockRepo := new(MockVenueRepository)
		svc := uc.NewService(mockRepo)
		handler := NewVenueHandler(svc)

		reqBody := map[string]interface{}{
			"name":      "Table 5",
			"capacity":  4,
			"can_merge": true,
			"zone":      "window",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/rooms/room-1/tables", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/rooms/:roomId/tables")
		c.SetParamNames("roomId")
		c.SetParamValues("room-1")

		expected := &venuepb.Table{Id: "table-new", Name: "Table 5", Capacity: 4}
		mockRepo.On("CreateTable", mock.Anything, mock.MatchedBy(func(r *venuepb.CreateTableRequest) bool {
			return r.RoomId == "room-1" && r.Name == "Table 5" && r.Capacity == 4
		})).Return(expected, nil)

		err := handler.CreateTable(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestVenueHandler_GetOpeningHours(t *testing.T) {
	e := echo.New()

	t.Run("successful get", func(t *testing.T) {
		mockRepo := new(MockVenueRepository)
		svc := uc.NewService(mockRepo)
		handler := NewVenueHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/venues/venue-1/schedule", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/venues/:venueId/schedule")
		c.SetParamNames("venueId")
		c.SetParamValues("venue-1")

		expected := &venuepb.OpeningHours{VenueId: "venue-1"}
		mockRepo.On("GetOpeningHours", mock.Anything, "venue-1").Return(expected, nil)

		err := handler.GetOpeningHours(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestVenueHandler_SetOpeningHours(t *testing.T) {
	e := echo.New()

	t.Run("successful set", func(t *testing.T) {
		mockRepo := new(MockVenueRepository)
		svc := uc.NewService(mockRepo)
		handler := NewVenueHandler(svc)

		reqBody := map[string]interface{}{
			"days": []map[string]interface{}{
				{"weekday": 1, "open_time": "09:00", "close_time": "22:00"},
			},
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/venues/venue-1/schedule", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/venues/:venueId/schedule")
		c.SetParamNames("venueId")
		c.SetParamValues("venue-1")

		expected := &venuepb.SetOpeningHoursResponse{Success: true}
		mockRepo.On("SetOpeningHours", mock.Anything, mock.MatchedBy(func(r *venuepb.SetOpeningHoursRequest) bool {
			return r.VenueId == "venue-1" && len(r.Days) == 1
		})).Return(expected, nil)

		err := handler.SetOpeningHours(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestVenueHandler_CheckAvailability(t *testing.T) {
	e := echo.New()

	t.Run("successful check", func(t *testing.T) {
		mockRepo := new(MockVenueRepository)
		svc := uc.NewService(mockRepo)
		handler := NewVenueHandler(svc)

		reqBody := map[string]interface{}{
			"venue_id": "venue-1",
			"slot": map[string]interface{}{
				"date":              "2025-11-12",
				"start_time":        "18:00",
				"duration_minutes":  120,
			},
			"party_size": 2,
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/availability/check", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		expected := &venuepb.CheckAvailabilityResponse{Tables: []*venuepb.TableAvailability{}}
		mockRepo.On("CheckAvailability", mock.Anything, mock.MatchedBy(func(r *venuepb.CheckAvailabilityRequest) bool {
			return r.VenueId == "venue-1" && r.PartySize == 2
		})).Return(expected, nil)

		err := handler.CheckAvailability(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestVenueHandler_GetRoom(t *testing.T) {
	e := echo.New()

	t.Run("successful get", func(t *testing.T) {
		mockRepo := new(MockVenueRepository)
		svc := uc.NewService(mockRepo)
		handler := NewVenueHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/rooms/room-1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/rooms/:id")
		c.SetParamNames("id")
		c.SetParamValues("room-1")

		expected := &venuepb.Room{Id: "room-1", Name: "Main Room"}
		mockRepo.On("GetRoom", mock.Anything, "room-1").Return(expected, nil)

		err := handler.GetRoom(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestVenueHandler_UpdateRoom(t *testing.T) {
	e := echo.New()

	t.Run("successful update", func(t *testing.T) {
		mockRepo := new(MockVenueRepository)
		svc := uc.NewService(mockRepo)
		handler := NewVenueHandler(svc)

		reqBody := map[string]interface{}{"name": "Updated Room"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPut, "/rooms/room-1", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/rooms/:id")
		c.SetParamNames("id")
		c.SetParamValues("room-1")

		expected := &venuepb.Room{Id: "room-1", Name: "Updated Room"}
		mockRepo.On("UpdateRoom", mock.Anything, mock.MatchedBy(func(r *venuepb.UpdateRoomRequest) bool {
			return r.Id == "room-1" && r.Name == "Updated Room"
		})).Return(expected, nil)

		err := handler.UpdateRoom(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestVenueHandler_DeleteRoom(t *testing.T) {
	e := echo.New()

	t.Run("successful delete", func(t *testing.T) {
		mockRepo := new(MockVenueRepository)
		svc := uc.NewService(mockRepo)
		handler := NewVenueHandler(svc)

		req := httptest.NewRequest(http.MethodDelete, "/rooms/room-1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/rooms/:id")
		c.SetParamNames("id")
		c.SetParamValues("room-1")

		mockRepo.On("DeleteRoom", mock.Anything, "room-1").Return(nil)

		err := handler.DeleteRoom(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestVenueHandler_GetTable(t *testing.T) {
	e := echo.New()

	t.Run("successful get", func(t *testing.T) {
		mockRepo := new(MockVenueRepository)
		svc := uc.NewService(mockRepo)
		handler := NewVenueHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/tables/table-1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/tables/:id")
		c.SetParamNames("id")
		c.SetParamValues("table-1")

		expected := &venuepb.Table{Id: "table-1", Name: "Table 1"}
		mockRepo.On("GetTable", mock.Anything, "table-1").Return(expected, nil)

		err := handler.GetTable(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestVenueHandler_UpdateTable(t *testing.T) {
	e := echo.New()

	t.Run("successful update", func(t *testing.T) {
		mockRepo := new(MockVenueRepository)
		svc := uc.NewService(mockRepo)
		handler := NewVenueHandler(svc)

		reqBody := map[string]interface{}{
			"name":      "Updated Table",
			"capacity":  6,
			"can_merge": false,
			"zone":      "outdoor",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPut, "/tables/table-1", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/tables/:id")
		c.SetParamNames("id")
		c.SetParamValues("table-1")

		expected := &venuepb.Table{Id: "table-1", Name: "Updated Table", Capacity: 6}
		mockRepo.On("UpdateTable", mock.Anything, mock.MatchedBy(func(r *venuepb.UpdateTableRequest) bool {
			return r.Id == "table-1" && r.Name == "Updated Table" && r.Capacity == 6
		})).Return(expected, nil)

		err := handler.UpdateTable(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestVenueHandler_DeleteTable(t *testing.T) {
	e := echo.New()

	t.Run("successful delete", func(t *testing.T) {
		mockRepo := new(MockVenueRepository)
		svc := uc.NewService(mockRepo)
		handler := NewVenueHandler(svc)

		req := httptest.NewRequest(http.MethodDelete, "/tables/table-1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/tables/:id")
		c.SetParamNames("id")
		c.SetParamValues("table-1")

		mockRepo.On("DeleteTable", mock.Anything, "table-1").Return(nil)

		err := handler.DeleteTable(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestVenueHandler_SetSpecialHours(t *testing.T) {
	e := echo.New()

	t.Run("successful set", func(t *testing.T) {
		mockRepo := new(MockVenueRepository)
		svc := uc.NewService(mockRepo)
		handler := NewVenueHandler(svc)

		reqBody := map[string]interface{}{
			"date":        "2025-12-25",
			"open_time":   "10:00",
			"close_time":  "20:00",
			"is_closed":   false,
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/venues/venue-1/special-hours", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/venues/:venueId/special-hours")
		c.SetParamNames("venueId")
		c.SetParamValues("venue-1")

		expected := &venuepb.SetSpecialHoursResponse{Success: true}
		mockRepo.On("SetSpecialHours", mock.Anything, mock.MatchedBy(func(r *venuepb.SetSpecialHoursRequest) bool {
			return r.VenueId == "venue-1" && r.Date == "2025-12-25" && !r.IsClosed
		})).Return(expected, nil)

		err := handler.SetSpecialHours(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

