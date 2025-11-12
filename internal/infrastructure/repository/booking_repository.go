package repository

import (
	"context"

	bookingpb "github.com/bookingcontrol/booker-contracts-go/booking"
	"github.com/bookingcontrol/booker-admin-gateway/internal/domain/repository"
)

type bookingRepository struct {
	client bookingpb.BookingServiceClient
}

func NewBookingRepository(client bookingpb.BookingServiceClient) repository.BookingRepository {
	return &bookingRepository{
		client: client,
	}
}

func (r *bookingRepository) ListBookings(ctx context.Context, req *bookingpb.ListBookingsRequest) (*bookingpb.ListBookingsResponse, error) {
	return r.client.ListBookings(ctx, req)
}

func (r *bookingRepository) GetBooking(ctx context.Context, id string) (*bookingpb.Booking, error) {
	return r.client.GetBooking(ctx, &bookingpb.GetBookingRequest{
		Id: id,
	})
}

func (r *bookingRepository) CreateBooking(ctx context.Context, req *bookingpb.CreateBookingRequest) (*bookingpb.Booking, error) {
	return r.client.CreateBooking(ctx, req)
}

func (r *bookingRepository) ConfirmBooking(ctx context.Context, id, adminID string) (*bookingpb.Booking, error) {
	return r.client.ConfirmBooking(ctx, &bookingpb.ConfirmBookingRequest{
		Id:      id,
		AdminId: adminID,
	})
}

func (r *bookingRepository) CancelBooking(ctx context.Context, id, adminID, reason string) (*bookingpb.Booking, error) {
	return r.client.CancelBooking(ctx, &bookingpb.CancelBookingRequest{
		Id:      id,
		AdminId: adminID,
		Reason:  reason,
	})
}

func (r *bookingRepository) MarkSeated(ctx context.Context, id, adminID string) (*bookingpb.Booking, error) {
	return r.client.MarkSeated(ctx, &bookingpb.MarkSeatedRequest{
		Id:      id,
		AdminId: adminID,
	})
}

func (r *bookingRepository) MarkFinished(ctx context.Context, id, adminID string) (*bookingpb.Booking, error) {
	return r.client.MarkFinished(ctx, &bookingpb.MarkFinishedRequest{
		Id:      id,
		AdminId: adminID,
	})
}

func (r *bookingRepository) MarkNoShow(ctx context.Context, id, adminID string) (*bookingpb.Booking, error) {
	return r.client.MarkNoShow(ctx, &bookingpb.MarkNoShowRequest{
		Id:      id,
		AdminId: adminID,
	})
}

