package usecase

import (
	"context"

	bookingpb "github.com/bookingcontrol/booker-contracts-go/booking"
	"github.com/bookingcontrol/booker-admin-gateway/internal/domain/repository"
)

type bookingUseCase struct {
	bookingRepo repository.BookingRepository
}

func NewBookingUseCase(bookingRepo repository.BookingRepository) *bookingUseCase {
	return &bookingUseCase{
		bookingRepo: bookingRepo,
	}
}

func (uc *bookingUseCase) ListBookings(ctx context.Context, req *bookingpb.ListBookingsRequest) (*bookingpb.ListBookingsResponse, error) {
	return uc.bookingRepo.ListBookings(ctx, req)
}

func (uc *bookingUseCase) GetBooking(ctx context.Context, id string) (*bookingpb.Booking, error) {
	return uc.bookingRepo.GetBooking(ctx, id)
}

func (uc *bookingUseCase) CreateBooking(ctx context.Context, req *bookingpb.CreateBookingRequest) (*bookingpb.Booking, error) {
	return uc.bookingRepo.CreateBooking(ctx, req)
}

func (uc *bookingUseCase) ConfirmBooking(ctx context.Context, id, adminID string) (*bookingpb.Booking, error) {
	return uc.bookingRepo.ConfirmBooking(ctx, id, adminID)
}

func (uc *bookingUseCase) CancelBooking(ctx context.Context, id, adminID, reason string) (*bookingpb.Booking, error) {
	return uc.bookingRepo.CancelBooking(ctx, id, adminID, reason)
}

func (uc *bookingUseCase) MarkSeated(ctx context.Context, id, adminID string) (*bookingpb.Booking, error) {
	return uc.bookingRepo.MarkSeated(ctx, id, adminID)
}

func (uc *bookingUseCase) MarkFinished(ctx context.Context, id, adminID string) (*bookingpb.Booking, error) {
	return uc.bookingRepo.MarkFinished(ctx, id, adminID)
}

func (uc *bookingUseCase) MarkNoShow(ctx context.Context, id, adminID string) (*bookingpb.Booking, error) {
	return uc.bookingRepo.MarkNoShow(ctx, id, adminID)
}

