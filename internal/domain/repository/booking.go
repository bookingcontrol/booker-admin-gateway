package repository

import (
	"context"

	bookingpb "github.com/bookingcontrol/booker-contracts-go/booking"
)

// BookingRepository defines interface for booking service operations
type BookingRepository interface {
	ListBookings(ctx context.Context, req *bookingpb.ListBookingsRequest) (*bookingpb.ListBookingsResponse, error)
	GetBooking(ctx context.Context, id string) (*bookingpb.Booking, error)
	CreateBooking(ctx context.Context, req *bookingpb.CreateBookingRequest) (*bookingpb.Booking, error)
	ConfirmBooking(ctx context.Context, id, adminID string) (*bookingpb.Booking, error)
	CancelBooking(ctx context.Context, id, adminID, reason string) (*bookingpb.Booking, error)
	MarkSeated(ctx context.Context, id, adminID string) (*bookingpb.Booking, error)
	MarkFinished(ctx context.Context, id, adminID string) (*bookingpb.Booking, error)
	MarkNoShow(ctx context.Context, id, adminID string) (*bookingpb.Booking, error)
}

