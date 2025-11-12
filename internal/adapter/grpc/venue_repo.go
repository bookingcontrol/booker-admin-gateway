package grpc

import (
	"context"

	venuepb "github.com/bookingcontrol/booker-contracts-go/venue"
	dom "github.com/bookingcontrol/booker-admin-gateway/internal/domain/venue"
)

type VenueRepo struct {
	client venuepb.VenueServiceClient
}

func NewVenueRepo(client venuepb.VenueServiceClient) dom.Repository {
	return &VenueRepo{
		client: client,
	}
}

func (r *VenueRepo) ListVenues(ctx context.Context, limit, offset int32) (*venuepb.ListVenuesResponse, error) {
	return r.client.ListVenues(ctx, &venuepb.ListVenuesRequest{
		Limit:  limit,
		Offset: offset,
	})
}

func (r *VenueRepo) GetVenue(ctx context.Context, id string) (*venuepb.Venue, error) {
	return r.client.GetVenue(ctx, &venuepb.GetVenueRequest{
		Id: id,
	})
}

func (r *VenueRepo) CreateVenue(ctx context.Context, req *venuepb.CreateVenueRequest) (*venuepb.Venue, error) {
	return r.client.CreateVenue(ctx, req)
}

func (r *VenueRepo) UpdateVenue(ctx context.Context, req *venuepb.UpdateVenueRequest) (*venuepb.Venue, error) {
	return r.client.UpdateVenue(ctx, req)
}

func (r *VenueRepo) DeleteVenue(ctx context.Context, id string) error {
	_, err := r.client.DeleteVenue(ctx, &venuepb.DeleteVenueRequest{
		Id: id,
	})
	return err
}

func (r *VenueRepo) ListRooms(ctx context.Context, venueID string, limit, offset int32) (*venuepb.ListRoomsResponse, error) {
	return r.client.ListRooms(ctx, &venuepb.ListRoomsRequest{
		VenueId: venueID,
		Limit:   limit,
		Offset:  offset,
	})
}

func (r *VenueRepo) GetRoom(ctx context.Context, id string) (*venuepb.Room, error) {
	return r.client.GetRoom(ctx, &venuepb.GetRoomRequest{
		Id: id,
	})
}

func (r *VenueRepo) CreateRoom(ctx context.Context, req *venuepb.CreateRoomRequest) (*venuepb.Room, error) {
	return r.client.CreateRoom(ctx, req)
}

func (r *VenueRepo) UpdateRoom(ctx context.Context, req *venuepb.UpdateRoomRequest) (*venuepb.Room, error) {
	return r.client.UpdateRoom(ctx, req)
}

func (r *VenueRepo) DeleteRoom(ctx context.Context, id string) error {
	_, err := r.client.DeleteRoom(ctx, &venuepb.DeleteRoomRequest{
		Id: id,
	})
	return err
}

func (r *VenueRepo) ListTables(ctx context.Context, roomID string, limit, offset int32) (*venuepb.ListTablesResponse, error) {
	return r.client.ListTables(ctx, &venuepb.ListTablesRequest{
		RoomId: roomID,
		Limit:  limit,
		Offset: offset,
	})
}

func (r *VenueRepo) GetTable(ctx context.Context, id string) (*venuepb.Table, error) {
	return r.client.GetTable(ctx, &venuepb.GetTableRequest{
		Id: id,
	})
}

func (r *VenueRepo) CreateTable(ctx context.Context, req *venuepb.CreateTableRequest) (*venuepb.Table, error) {
	return r.client.CreateTable(ctx, req)
}

func (r *VenueRepo) UpdateTable(ctx context.Context, req *venuepb.UpdateTableRequest) (*venuepb.Table, error) {
	return r.client.UpdateTable(ctx, req)
}

func (r *VenueRepo) DeleteTable(ctx context.Context, id string) error {
	_, err := r.client.DeleteTable(ctx, &venuepb.DeleteTableRequest{
		Id: id,
	})
	return err
}

func (r *VenueRepo) GetOpeningHours(ctx context.Context, venueID string) (*venuepb.OpeningHours, error) {
	return r.client.GetOpeningHours(ctx, &venuepb.GetOpeningHoursRequest{
		VenueId: venueID,
	})
}

func (r *VenueRepo) SetOpeningHours(ctx context.Context, req *venuepb.SetOpeningHoursRequest) (*venuepb.SetOpeningHoursResponse, error) {
	return r.client.SetOpeningHours(ctx, req)
}

func (r *VenueRepo) SetSpecialHours(ctx context.Context, req *venuepb.SetSpecialHoursRequest) (*venuepb.SetSpecialHoursResponse, error) {
	return r.client.SetSpecialHours(ctx, req)
}

func (r *VenueRepo) CheckAvailability(ctx context.Context, req *venuepb.CheckAvailabilityRequest) (*venuepb.CheckAvailabilityResponse, error) {
	return r.client.CheckAvailability(ctx, req)
}

