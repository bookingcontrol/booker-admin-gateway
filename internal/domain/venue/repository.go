package venue

import (
	"context"

	venuepb "github.com/bookingcontrol/booker-contracts-go/venue"
)

// Repository defines interface for venue service operations
type Repository interface {
	ListVenues(ctx context.Context, limit, offset int32) (*venuepb.ListVenuesResponse, error)
	GetVenue(ctx context.Context, id string) (*venuepb.Venue, error)
	CreateVenue(ctx context.Context, req *venuepb.CreateVenueRequest) (*venuepb.Venue, error)
	UpdateVenue(ctx context.Context, req *venuepb.UpdateVenueRequest) (*venuepb.Venue, error)
	DeleteVenue(ctx context.Context, id string) error
	ListRooms(ctx context.Context, venueID string, limit, offset int32) (*venuepb.ListRoomsResponse, error)
	GetRoom(ctx context.Context, id string) (*venuepb.Room, error)
	CreateRoom(ctx context.Context, req *venuepb.CreateRoomRequest) (*venuepb.Room, error)
	UpdateRoom(ctx context.Context, req *venuepb.UpdateRoomRequest) (*venuepb.Room, error)
	DeleteRoom(ctx context.Context, id string) error
	ListTables(ctx context.Context, roomID string, limit, offset int32) (*venuepb.ListTablesResponse, error)
	GetTable(ctx context.Context, id string) (*venuepb.Table, error)
	CreateTable(ctx context.Context, req *venuepb.CreateTableRequest) (*venuepb.Table, error)
	UpdateTable(ctx context.Context, req *venuepb.UpdateTableRequest) (*venuepb.Table, error)
	DeleteTable(ctx context.Context, id string) error
	GetOpeningHours(ctx context.Context, venueID string) (*venuepb.OpeningHours, error)
	SetOpeningHours(ctx context.Context, req *venuepb.SetOpeningHoursRequest) (*venuepb.SetOpeningHoursResponse, error)
	SetSpecialHours(ctx context.Context, req *venuepb.SetSpecialHoursRequest) (*venuepb.SetSpecialHoursResponse, error)
	CheckAvailability(ctx context.Context, req *venuepb.CheckAvailabilityRequest) (*venuepb.CheckAvailabilityResponse, error)
}

