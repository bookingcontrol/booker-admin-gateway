package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/protobuf/encoding/protojson"

	bookingpb "github.com/bookingcontrol/booker-contracts-go/booking"
	commonpb "github.com/bookingcontrol/booker-contracts-go/common"
	venuepb "github.com/bookingcontrol/booker-contracts-go/venue"
)

// Booking handlers
func (h *Handler) ListBookings(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 50
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	resp, err := h.bookingUseCase.ListBookings(c.Request().Context(), &bookingpb.ListBookingsRequest{
		VenueId: c.QueryParam("venue_id"),
		Date:    c.QueryParam("date"),
		Status:  c.QueryParam("status"),
		TableId: c.QueryParam("table_id"),
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetBooking(c echo.Context) error {
	resp, err := h.bookingUseCase.GetBooking(c.Request().Context(), c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) CreateBooking(c echo.Context) error {
	var req struct {
		VenueID string `json:"venue_id"`
		Table   struct {
			VenueID string `json:"venue_id"`
			RoomID  string `json:"room_id"`
			TableID string `json:"table_id"`
		} `json:"table"`
		Slot struct {
			Date            string `json:"date"`
			StartTime       string `json:"start_time"`
			DurationMinutes int32  `json:"duration_minutes"`
		} `json:"slot"`
		PartySize      int32  `json:"party_size"`
		CustomerName   string `json:"customer_name"`
		CustomerPhone  string `json:"customer_phone"`
		Comment        string `json:"comment"`
		IdempotencyKey string `json:"idempotency_key"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	adminID := c.Get("admin_id").(string)

	resp, err := h.bookingUseCase.CreateBooking(c.Request().Context(), &bookingpb.CreateBookingRequest{
		VenueId: req.VenueID,
		Table: &commonpb.TableRef{
			VenueId: req.Table.VenueID,
			RoomId:  req.Table.RoomID,
			TableId: req.Table.TableID,
		},
		Slot: &commonpb.Slot{
			Date:            req.Slot.Date,
			StartTime:       req.Slot.StartTime,
			DurationMinutes: req.Slot.DurationMinutes,
		},
		PartySize:      req.PartySize,
		CustomerName:   req.CustomerName,
		CustomerPhone:  req.CustomerPhone,
		Comment:        req.Comment,
		AdminId:        adminID,
		IdempotencyKey: req.IdempotencyKey,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h *Handler) ConfirmBooking(c echo.Context) error {
	adminID := c.Get("admin_id").(string)

	resp, err := h.bookingUseCase.ConfirmBooking(c.Request().Context(), c.Param("id"), adminID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) CancelBooking(c echo.Context) error {
	var req struct {
		Reason string `json:"reason"`
	}
	c.Bind(&req)

	adminID := c.Get("admin_id").(string)

	resp, err := h.bookingUseCase.CancelBooking(c.Request().Context(), c.Param("id"), adminID, req.Reason)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) MarkSeated(c echo.Context) error {
	adminID := c.Get("admin_id").(string)

	resp, err := h.bookingUseCase.MarkSeated(c.Request().Context(), c.Param("id"), adminID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) MarkFinished(c echo.Context) error {
	adminID := c.Get("admin_id").(string)

	resp, err := h.bookingUseCase.MarkFinished(c.Request().Context(), c.Param("id"), adminID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) MarkNoShow(c echo.Context) error {
	adminID := c.Get("admin_id").(string)

	resp, err := h.bookingUseCase.MarkNoShow(c.Request().Context(), c.Param("id"), adminID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) CheckAvailability(c echo.Context) error {
	var req struct {
		VenueID string `json:"venue_id"`
		Slot    struct {
			Date            string `json:"date"`
			StartTime       string `json:"start_time"`
			DurationMinutes int32  `json:"duration_minutes"`
		} `json:"slot"`
		PartySize int32 `json:"party_size"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	resp, err := h.venueUseCase.CheckAvailability(c.Request().Context(), &venuepb.CheckAvailabilityRequest{
		VenueId: req.VenueID,
		Slot: &commonpb.Slot{
			Date:            req.Slot.Date,
			StartTime:       req.Slot.StartTime,
			DurationMinutes: req.Slot.DurationMinutes,
		},
		PartySize: req.PartySize,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Use protojson to properly serialize protobuf message to JSON
	// This ensures all fields including merged_with_table are included
	marshaler := protojson.MarshalOptions{
		EmitUnpopulated: true,
		UseProtoNames:   true, // Use snake_case field names
	}
	jsonBytes, err := marshaler.Marshal(resp)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.Blob(http.StatusOK, echo.MIMEApplicationJSON, jsonBytes)
}

func (h *Handler) WebSocket(c echo.Context) error {
	// TODO: Implement WebSocket for live updates
	return c.String(http.StatusNotImplemented, "WebSocket not implemented yet")
}

// Metrics endpoint
func (h *Handler) Metrics(c echo.Context) error {
	promhttp.Handler().ServeHTTP(c.Response(), c.Request())
	return nil
}

