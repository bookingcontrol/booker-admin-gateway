package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	bookingpb "github.com/bookingcontrol/booker-contracts-go/booking"
	commonpb "github.com/bookingcontrol/booker-contracts-go/common"
	uc "github.com/bookingcontrol/booker-admin-gateway/internal/usecase/booking"
)

type BookingHandler struct {
	svc *uc.Service
}

func NewBookingHandler(svc *uc.Service) *BookingHandler {
	return &BookingHandler{svc: svc}
}

func (h *BookingHandler) ListBookings(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 50
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	resp, err := h.svc.ListBookings(c.Request().Context(), &bookingpb.ListBookingsRequest{
		VenueId: c.QueryParam("venue_id"), Date: c.QueryParam("date"), Status: c.QueryParam("status"),
		TableId: c.QueryParam("table_id"), Limit: int32(limit), Offset: int32(offset),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *BookingHandler) GetBooking(c echo.Context) error {
	resp, err := h.svc.GetBooking(c.Request().Context(), c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *BookingHandler) CreateBooking(c echo.Context) error {
	var req struct {
		VenueID        string
		Table          struct{ VenueID, RoomID, TableID string }
		Slot           struct{ Date, StartTime string; DurationMinutes int32 }
		PartySize      int32
		CustomerName   string
		CustomerPhone  string
		Comment        string
		IdempotencyKey string
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	adminID := c.Get("admin_id").(string)
	resp, err := h.svc.CreateBooking(c.Request().Context(), &bookingpb.CreateBookingRequest{
		VenueId: req.VenueID,
		Table:   &commonpb.TableRef{VenueId: req.Table.VenueID, RoomId: req.Table.RoomID, TableId: req.Table.TableID},
		Slot:    &commonpb.Slot{Date: req.Slot.Date, StartTime: req.Slot.StartTime, DurationMinutes: req.Slot.DurationMinutes},
		PartySize: req.PartySize, CustomerName: req.CustomerName, CustomerPhone: req.CustomerPhone,
		Comment: req.Comment, AdminId: adminID, IdempotencyKey: req.IdempotencyKey,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, resp)
}

func (h *BookingHandler) ConfirmBooking(c echo.Context) error {
	adminID := c.Get("admin_id").(string)
	resp, err := h.svc.ConfirmBooking(c.Request().Context(), c.Param("id"), adminID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *BookingHandler) CancelBooking(c echo.Context) error {
	var req struct{ Reason string }
	c.Bind(&req)
	adminID := c.Get("admin_id").(string)
	resp, err := h.svc.CancelBooking(c.Request().Context(), c.Param("id"), adminID, req.Reason)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *BookingHandler) MarkSeated(c echo.Context) error {
	adminID := c.Get("admin_id").(string)
	resp, err := h.svc.MarkSeated(c.Request().Context(), c.Param("id"), adminID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *BookingHandler) MarkFinished(c echo.Context) error {
	adminID := c.Get("admin_id").(string)
	resp, err := h.svc.MarkFinished(c.Request().Context(), c.Param("id"), adminID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *BookingHandler) MarkNoShow(c echo.Context) error {
	adminID := c.Get("admin_id").(string)
	resp, err := h.svc.MarkNoShow(c.Request().Context(), c.Param("id"), adminID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *BookingHandler) WebSocket(c echo.Context) error {
	return c.String(http.StatusNotImplemented, "WebSocket not implemented yet")
}

func (h *BookingHandler) Metrics(c echo.Context) error {
	promhttp.Handler().ServeHTTP(c.Response(), c.Request())
	return nil
}

