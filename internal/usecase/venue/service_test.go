package venue

import (
	"context"
	"errors"
	"testing"

	venuepb "github.com/bookingcontrol/booker-contracts-go/venue"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockVenueRepository is a mock implementation of venue repository
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

func TestService_ListVenues(t *testing.T) {
	t.Run("successful list", func(t *testing.T) {
		mockRepo := new(MockVenueRepository)
		service := NewService(mockRepo)

		expected := &venuepb.ListVenuesResponse{
			Venues: []*venuepb.Venue{
				{Id: "venue-1", Name: "Test Venue"},
			},
			Total: 1,
		}
		mockRepo.On("ListVenues", mock.Anything, int32(50), int32(0)).Return(expected, nil)

		result, err := service.ListVenues(context.Background(), 50, 0)

		require.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockVenueRepository)
		service := NewService(mockRepo)

		mockRepo.On("ListVenues", mock.Anything, int32(50), int32(0)).Return(nil, errors.New("db error"))

		_, err := service.ListVenues(context.Background(), 50, 0)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestService_GetVenue(t *testing.T) {
	t.Run("successful get", func(t *testing.T) {
		mockRepo := new(MockVenueRepository)
		service := NewService(mockRepo)

		expected := &venuepb.Venue{Id: "venue-1", Name: "Test Venue"}
		mockRepo.On("GetVenue", mock.Anything, "venue-1").Return(expected, nil)

		result, err := service.GetVenue(context.Background(), "venue-1")

		require.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("venue not found", func(t *testing.T) {
		mockRepo := new(MockVenueRepository)
		service := NewService(mockRepo)

		mockRepo.On("GetVenue", mock.Anything, "nonexistent").Return(nil, errors.New("not found"))

		_, err := service.GetVenue(context.Background(), "nonexistent")

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestService_CreateVenue(t *testing.T) {
	t.Run("successful create", func(t *testing.T) {
		mockRepo := new(MockVenueRepository)
		service := NewService(mockRepo)

		req := &venuepb.CreateVenueRequest{
			Name:     "New Venue",
			Timezone: "UTC",
			Address:  "123 Main St",
			Phone:    "+1234567890",
			Email:    "venue@example.com",
		}
		expected := &venuepb.Venue{
			Id:       "venue-new",
			Name:     "New Venue",
			Timezone: "UTC",
			Address:  "123 Main St",
			Phone:    "+1234567890",
			Email:    "venue@example.com",
		}
		mockRepo.On("CreateVenue", mock.Anything, req).Return(expected, nil)

		result, err := service.CreateVenue(context.Background(), req)

		require.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestService_DeleteVenue(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		mockRepo := new(MockVenueRepository)
		service := NewService(mockRepo)

		mockRepo.On("DeleteVenue", mock.Anything, "venue-1").Return(nil)

		err := service.DeleteVenue(context.Background(), "venue-1")

		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockVenueRepository)
		service := NewService(mockRepo)

		mockRepo.On("DeleteVenue", mock.Anything, "venue-1").Return(errors.New("db error"))

		err := service.DeleteVenue(context.Background(), "venue-1")

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
