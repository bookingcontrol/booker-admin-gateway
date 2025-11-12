package booking

import (
	"context"

	bookingpb "github.com/bookingcontrol/booker-contracts-go/booking"
	dom "github.com/bookingcontrol/booker-admin-gateway/internal/domain/booking"
)

type Service struct {
	repo dom.Repository
}

func NewService(repo dom.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) ListBookings(ctx context.Context, req *bookingpb.ListBookingsRequest) (*bookingpb.ListBookingsResponse, error) {
	return s.repo.ListBookings(ctx, req)
}

func (s *Service) GetBooking(ctx context.Context, id string) (*bookingpb.Booking, error) {
	return s.repo.GetBooking(ctx, id)
}

func (s *Service) CreateBooking(ctx context.Context, req *bookingpb.CreateBookingRequest) (*bookingpb.Booking, error) {
	return s.repo.CreateBooking(ctx, req)
}

func (s *Service) ConfirmBooking(ctx context.Context, id, adminID string) (*bookingpb.Booking, error) {
	return s.repo.ConfirmBooking(ctx, id, adminID)
}

func (s *Service) CancelBooking(ctx context.Context, id, adminID, reason string) (*bookingpb.Booking, error) {
	return s.repo.CancelBooking(ctx, id, adminID, reason)
}

func (s *Service) MarkSeated(ctx context.Context, id, adminID string) (*bookingpb.Booking, error) {
	return s.repo.MarkSeated(ctx, id, adminID)
}

func (s *Service) MarkFinished(ctx context.Context, id, adminID string) (*bookingpb.Booking, error) {
	return s.repo.MarkFinished(ctx, id, adminID)
}

func (s *Service) MarkNoShow(ctx context.Context, id, adminID string) (*bookingpb.Booking, error) {
	return s.repo.MarkNoShow(ctx, id, adminID)
}

