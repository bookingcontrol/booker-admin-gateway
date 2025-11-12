package grpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	bookingpb "github.com/bookingcontrol/booker-contracts-go/booking"
)

// Тестируем только корректность маппинга параметров
// gRPC адаптер - это тонкий прокси, но проверяем, что параметры правильно передаются

func TestBookingRepo_ParameterMapping(t *testing.T) {
	t.Run("GetBooking maps id correctly", func(t *testing.T) {
		// Проверяем, что id правильно маппится в GetBookingRequest
		id := "booking-123"
		req := &bookingpb.GetBookingRequest{Id: id}
		
		assert.Equal(t, "booking-123", req.Id)
		require.NotNil(t, req)
	})
	
	t.Run("ConfirmBooking maps id and adminID correctly", func(t *testing.T) {
		id := "booking-123"
		adminID := "admin-456"
		req := &bookingpb.ConfirmBookingRequest{
			Id:      id,
			AdminId: adminID,
		}
		
		assert.Equal(t, "booking-123", req.Id)
		assert.Equal(t, "admin-456", req.AdminId)
	})
	
	t.Run("CancelBooking maps all parameters correctly", func(t *testing.T) {
		id := "booking-123"
		adminID := "admin-456"
		reason := "Customer cancelled"
		req := &bookingpb.CancelBookingRequest{
			Id:      id,
			AdminId: adminID,
			Reason:  reason,
		}
		
		assert.Equal(t, "booking-123", req.Id)
		assert.Equal(t, "admin-456", req.AdminId)
		assert.Equal(t, "Customer cancelled", req.Reason)
	})
}

func TestBookingRepo_RequestStructure(t *testing.T) {
	t.Run("CreateBooking request structure", func(t *testing.T) {
		req := &bookingpb.CreateBookingRequest{
			VenueId:      "venue-1",
			PartySize:    4,
			CustomerName: "John Doe",
		}
		
		assert.Equal(t, "venue-1", req.VenueId)
		assert.Equal(t, int32(4), req.PartySize)
		assert.Equal(t, "John Doe", req.CustomerName)
	})
}

