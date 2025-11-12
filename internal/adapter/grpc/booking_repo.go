package grpc

import (
	"context"

	bookingpb "github.com/bookingcontrol/booker-contracts-go/booking"
	dom "github.com/bookingcontrol/booker-admin-gateway/internal/domain/booking"
)

type BookingRepo struct {
	client bookingpb.BookingServiceClient
}

func NewBookingRepo(client bookingpb.BookingServiceClient) dom.Repository {
	return &BookingRepo{
		client: client,
	}
}

func (r *BookingRepo) ListBookings(ctx context.Context, req *bookingpb.ListBookingsRequest) (*bookingpb.ListBookingsResponse, error) {
	return r.client.ListBookings(ctx, req)
}

func (r *BookingRepo) GetBooking(ctx context.Context, id string) (*bookingpb.Booking, error) {
	return r.client.GetBooking(ctx, &bookingpb.GetBookingRequest{
		Id: id,
	})
}

func (r *BookingRepo) CreateBooking(ctx context.Context, req *bookingpb.CreateBookingRequest) (*bookingpb.Booking, error) {
	return r.client.CreateBooking(ctx, req)
}

func (r *BookingRepo) ConfirmBooking(ctx context.Context, id, adminID string) (*bookingpb.Booking, error) {
	return r.client.ConfirmBooking(ctx, &bookingpb.ConfirmBookingRequest{
		Id:      id,
		AdminId: adminID,
	})
}

func (r *BookingRepo) CancelBooking(ctx context.Context, id, adminID, reason string) (*bookingpb.Booking, error) {
	return r.client.CancelBooking(ctx, &bookingpb.CancelBookingRequest{
		Id:      id,
		AdminId: adminID,
		Reason:  reason,
	})
}

func (r *BookingRepo) MarkSeated(ctx context.Context, id, adminID string) (*bookingpb.Booking, error) {
	return r.client.MarkSeated(ctx, &bookingpb.MarkSeatedRequest{
		Id:      id,
		AdminId: adminID,
	})
}

func (r *BookingRepo) MarkFinished(ctx context.Context, id, adminID string) (*bookingpb.Booking, error) {
	return r.client.MarkFinished(ctx, &bookingpb.MarkFinishedRequest{
		Id:      id,
		AdminId: adminID,
	})
}

func (r *BookingRepo) MarkNoShow(ctx context.Context, id, adminID string) (*bookingpb.Booking, error) {
	return r.client.MarkNoShow(ctx, &bookingpb.MarkNoShowRequest{
		Id:      id,
		AdminId: adminID,
	})
}

