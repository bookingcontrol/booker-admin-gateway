package repository

import (
	"context"

	venuepb "github.com/bookingcontrol/booker-contracts-go/venue"
	"github.com/bookingcontrol/booker-admin-gateway/internal/domain/repository"
)

type venueRepository struct {
	client venuepb.VenueServiceClient
}

func NewVenueRepository(client venuepb.VenueServiceClient) repository.VenueRepository {
	return &venueRepository{
		client: client,
	}
}

func (r *venueRepository) ListVenues(ctx context.Context, limit, offset int32) (*venuepb.ListVenuesResponse, error) {
	return r.client.ListVenues(ctx, &venuepb.ListVenuesRequest{
		Limit:  limit,
		Offset: offset,
	})
}

func (r *venueRepository) GetVenue(ctx context.Context, id string) (*venuepb.Venue, error) {
	return r.client.GetVenue(ctx, &venuepb.GetVenueRequest{
		Id: id,
	})
}

func (r *venueRepository) CreateVenue(ctx context.Context, req *venuepb.CreateVenueRequest) (*venuepb.Venue, error) {
	return r.client.CreateVenue(ctx, req)
}

func (r *venueRepository) UpdateVenue(ctx context.Context, req *venuepb.UpdateVenueRequest) (*venuepb.Venue, error) {
	return r.client.UpdateVenue(ctx, req)
}

func (r *venueRepository) DeleteVenue(ctx context.Context, id string) error {
	_, err := r.client.DeleteVenue(ctx, &venuepb.DeleteVenueRequest{
		Id: id,
	})
	return err
}

func (r *venueRepository) ListRooms(ctx context.Context, venueID string, limit, offset int32) (*venuepb.ListRoomsResponse, error) {
	return r.client.ListRooms(ctx, &venuepb.ListRoomsRequest{
		VenueId: venueID,
		Limit:   limit,
		Offset:  offset,
	})
}

func (r *venueRepository) GetRoom(ctx context.Context, id string) (*venuepb.Room, error) {
	return r.client.GetRoom(ctx, &venuepb.GetRoomRequest{
		Id: id,
	})
}

func (r *venueRepository) CreateRoom(ctx context.Context, req *venuepb.CreateRoomRequest) (*venuepb.Room, error) {
	return r.client.CreateRoom(ctx, req)
}

func (r *venueRepository) UpdateRoom(ctx context.Context, req *venuepb.UpdateRoomRequest) (*venuepb.Room, error) {
	return r.client.UpdateRoom(ctx, req)
}

func (r *venueRepository) DeleteRoom(ctx context.Context, id string) error {
	_, err := r.client.DeleteRoom(ctx, &venuepb.DeleteRoomRequest{
		Id: id,
	})
	return err
}

func (r *venueRepository) ListTables(ctx context.Context, roomID string, limit, offset int32) (*venuepb.ListTablesResponse, error) {
	return r.client.ListTables(ctx, &venuepb.ListTablesRequest{
		RoomId: roomID,
		Limit:  limit,
		Offset: offset,
	})
}

func (r *venueRepository) GetTable(ctx context.Context, id string) (*venuepb.Table, error) {
	return r.client.GetTable(ctx, &venuepb.GetTableRequest{
		Id: id,
	})
}

func (r *venueRepository) CreateTable(ctx context.Context, req *venuepb.CreateTableRequest) (*venuepb.Table, error) {
	return r.client.CreateTable(ctx, req)
}

func (r *venueRepository) UpdateTable(ctx context.Context, req *venuepb.UpdateTableRequest) (*venuepb.Table, error) {
	return r.client.UpdateTable(ctx, req)
}

func (r *venueRepository) DeleteTable(ctx context.Context, id string) error {
	_, err := r.client.DeleteTable(ctx, &venuepb.DeleteTableRequest{
		Id: id,
	})
	return err
}

func (r *venueRepository) GetOpeningHours(ctx context.Context, venueID string) (*venuepb.OpeningHours, error) {
	return r.client.GetOpeningHours(ctx, &venuepb.GetOpeningHoursRequest{
		VenueId: venueID,
	})
}

func (r *venueRepository) SetOpeningHours(ctx context.Context, req *venuepb.SetOpeningHoursRequest) (*venuepb.SetOpeningHoursResponse, error) {
	return r.client.SetOpeningHours(ctx, req)
}

func (r *venueRepository) SetSpecialHours(ctx context.Context, req *venuepb.SetSpecialHoursRequest) (*venuepb.SetSpecialHoursResponse, error) {
	return r.client.SetSpecialHours(ctx, req)
}

func (r *venueRepository) CheckAvailability(ctx context.Context, req *venuepb.CheckAvailabilityRequest) (*venuepb.CheckAvailabilityResponse, error) {
	return r.client.CheckAvailability(ctx, req)
}

