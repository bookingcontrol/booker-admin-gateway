package venue

import (
	"context"

	venuepb "github.com/bookingcontrol/booker-contracts-go/venue"
	dom "github.com/bookingcontrol/booker-admin-gateway/internal/domain/venue"
)

type Service struct {
	repo dom.Repository
}

func NewService(repo dom.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) ListVenues(ctx context.Context, limit, offset int32) (*venuepb.ListVenuesResponse, error) {
	return s.repo.ListVenues(ctx, limit, offset)
}

func (s *Service) GetVenue(ctx context.Context, id string) (*venuepb.Venue, error) {
	return s.repo.GetVenue(ctx, id)
}

func (s *Service) CreateVenue(ctx context.Context, req *venuepb.CreateVenueRequest) (*venuepb.Venue, error) {
	return s.repo.CreateVenue(ctx, req)
}

func (s *Service) UpdateVenue(ctx context.Context, req *venuepb.UpdateVenueRequest) (*venuepb.Venue, error) {
	return s.repo.UpdateVenue(ctx, req)
}

func (s *Service) DeleteVenue(ctx context.Context, id string) error {
	return s.repo.DeleteVenue(ctx, id)
}

func (s *Service) ListRooms(ctx context.Context, venueID string, limit, offset int32) (*venuepb.ListRoomsResponse, error) {
	return s.repo.ListRooms(ctx, venueID, limit, offset)
}

func (s *Service) GetRoom(ctx context.Context, id string) (*venuepb.Room, error) {
	return s.repo.GetRoom(ctx, id)
}

func (s *Service) CreateRoom(ctx context.Context, req *venuepb.CreateRoomRequest) (*venuepb.Room, error) {
	return s.repo.CreateRoom(ctx, req)
}

func (s *Service) UpdateRoom(ctx context.Context, req *venuepb.UpdateRoomRequest) (*venuepb.Room, error) {
	return s.repo.UpdateRoom(ctx, req)
}

func (s *Service) DeleteRoom(ctx context.Context, id string) error {
	return s.repo.DeleteRoom(ctx, id)
}

func (s *Service) ListTables(ctx context.Context, roomID string, limit, offset int32) (*venuepb.ListTablesResponse, error) {
	return s.repo.ListTables(ctx, roomID, limit, offset)
}

func (s *Service) GetTable(ctx context.Context, id string) (*venuepb.Table, error) {
	return s.repo.GetTable(ctx, id)
}

func (s *Service) CreateTable(ctx context.Context, req *venuepb.CreateTableRequest) (*venuepb.Table, error) {
	return s.repo.CreateTable(ctx, req)
}

func (s *Service) UpdateTable(ctx context.Context, req *venuepb.UpdateTableRequest) (*venuepb.Table, error) {
	return s.repo.UpdateTable(ctx, req)
}

func (s *Service) DeleteTable(ctx context.Context, id string) error {
	return s.repo.DeleteTable(ctx, id)
}

func (s *Service) GetOpeningHours(ctx context.Context, venueID string) (*venuepb.OpeningHours, error) {
	return s.repo.GetOpeningHours(ctx, venueID)
}

func (s *Service) SetOpeningHours(ctx context.Context, req *venuepb.SetOpeningHoursRequest) (*venuepb.SetOpeningHoursResponse, error) {
	return s.repo.SetOpeningHours(ctx, req)
}

func (s *Service) SetSpecialHours(ctx context.Context, req *venuepb.SetSpecialHoursRequest) (*venuepb.SetSpecialHoursResponse, error) {
	return s.repo.SetSpecialHours(ctx, req)
}

func (s *Service) CheckAvailability(ctx context.Context, req *venuepb.CheckAvailabilityRequest) (*venuepb.CheckAvailabilityResponse, error) {
	return s.repo.CheckAvailability(ctx, req)
}

