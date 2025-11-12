package usecase

import (
	"context"

	venuepb "github.com/bookingcontrol/booker-contracts-go/venue"
	"github.com/bookingcontrol/booker-admin-gateway/internal/domain/repository"
)

type venueUseCase struct {
	venueRepo repository.VenueRepository
}

func NewVenueUseCase(venueRepo repository.VenueRepository) *venueUseCase {
	return &venueUseCase{
		venueRepo: venueRepo,
	}
}

func (uc *venueUseCase) ListVenues(ctx context.Context, limit, offset int32) (*venuepb.ListVenuesResponse, error) {
	return uc.venueRepo.ListVenues(ctx, limit, offset)
}

func (uc *venueUseCase) GetVenue(ctx context.Context, id string) (*venuepb.Venue, error) {
	return uc.venueRepo.GetVenue(ctx, id)
}

func (uc *venueUseCase) CreateVenue(ctx context.Context, req *venuepb.CreateVenueRequest) (*venuepb.Venue, error) {
	return uc.venueRepo.CreateVenue(ctx, req)
}

func (uc *venueUseCase) UpdateVenue(ctx context.Context, req *venuepb.UpdateVenueRequest) (*venuepb.Venue, error) {
	return uc.venueRepo.UpdateVenue(ctx, req)
}

func (uc *venueUseCase) DeleteVenue(ctx context.Context, id string) error {
	return uc.venueRepo.DeleteVenue(ctx, id)
}

func (uc *venueUseCase) ListRooms(ctx context.Context, venueID string, limit, offset int32) (*venuepb.ListRoomsResponse, error) {
	return uc.venueRepo.ListRooms(ctx, venueID, limit, offset)
}

func (uc *venueUseCase) GetRoom(ctx context.Context, id string) (*venuepb.Room, error) {
	return uc.venueRepo.GetRoom(ctx, id)
}

func (uc *venueUseCase) CreateRoom(ctx context.Context, req *venuepb.CreateRoomRequest) (*venuepb.Room, error) {
	return uc.venueRepo.CreateRoom(ctx, req)
}

func (uc *venueUseCase) UpdateRoom(ctx context.Context, req *venuepb.UpdateRoomRequest) (*venuepb.Room, error) {
	return uc.venueRepo.UpdateRoom(ctx, req)
}

func (uc *venueUseCase) DeleteRoom(ctx context.Context, id string) error {
	return uc.venueRepo.DeleteRoom(ctx, id)
}

func (uc *venueUseCase) ListTables(ctx context.Context, roomID string, limit, offset int32) (*venuepb.ListTablesResponse, error) {
	return uc.venueRepo.ListTables(ctx, roomID, limit, offset)
}

func (uc *venueUseCase) GetTable(ctx context.Context, id string) (*venuepb.Table, error) {
	return uc.venueRepo.GetTable(ctx, id)
}

func (uc *venueUseCase) CreateTable(ctx context.Context, req *venuepb.CreateTableRequest) (*venuepb.Table, error) {
	return uc.venueRepo.CreateTable(ctx, req)
}

func (uc *venueUseCase) UpdateTable(ctx context.Context, req *venuepb.UpdateTableRequest) (*venuepb.Table, error) {
	return uc.venueRepo.UpdateTable(ctx, req)
}

func (uc *venueUseCase) DeleteTable(ctx context.Context, id string) error {
	return uc.venueRepo.DeleteTable(ctx, id)
}

func (uc *venueUseCase) GetOpeningHours(ctx context.Context, venueID string) (*venuepb.OpeningHours, error) {
	return uc.venueRepo.GetOpeningHours(ctx, venueID)
}

func (uc *venueUseCase) SetOpeningHours(ctx context.Context, req *venuepb.SetOpeningHoursRequest) (*venuepb.SetOpeningHoursResponse, error) {
	return uc.venueRepo.SetOpeningHours(ctx, req)
}

func (uc *venueUseCase) SetSpecialHours(ctx context.Context, req *venuepb.SetSpecialHoursRequest) (*venuepb.SetSpecialHoursResponse, error) {
	return uc.venueRepo.SetSpecialHours(ctx, req)
}

func (uc *venueUseCase) CheckAvailability(ctx context.Context, req *venuepb.CheckAvailabilityRequest) (*venuepb.CheckAvailabilityResponse, error) {
	return uc.venueRepo.CheckAvailability(ctx, req)
}

