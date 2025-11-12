package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
	venuepb "github.com/bookingcontrol/booker-contracts-go/venue"
	commonpb "github.com/bookingcontrol/booker-contracts-go/common"
	uc "github.com/bookingcontrol/booker-admin-gateway/internal/usecase/venue"
)

type VenueHandler struct {
	svc *uc.Service
}

func NewVenueHandler(svc *uc.Service) *VenueHandler {
	return &VenueHandler{svc: svc}
}

func (h *VenueHandler) ListVenues(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 50
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	resp, err := h.svc.ListVenues(c.Request().Context(), int32(limit), int32(offset))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *VenueHandler) GetVenue(c echo.Context) error {
	resp, err := h.svc.GetVenue(c.Request().Context(), c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *VenueHandler) CreateVenue(c echo.Context) error {
	var req struct {
		Name     string `json:"name"`
		Timezone string `json:"timezone"`
		Address  string `json:"address"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
	}
	if err := c.Bind(&req); err != nil {
		log.Warn().Err(err).Msg("Failed to bind CreateVenue request")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	log.Info().Str("name", req.Name).Msg("Creating venue")
	resp, err := h.svc.CreateVenue(c.Request().Context(), &venuepb.CreateVenueRequest{
		Name: req.Name, Timezone: req.Timezone, Address: req.Address, Phone: req.Phone, Email: req.Email,
	})
	if err != nil {
		log.Error().Err(err).Str("name", req.Name).Msg("Failed to create venue")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	log.Info().Str("venue_id", resp.Id).Str("name", resp.Name).Msg("Venue created successfully")
	return c.JSON(http.StatusCreated, resp)
}

func (h *VenueHandler) UpdateVenue(c echo.Context) error {
	var req struct {
		Name, Address, Phone, Email string
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	resp, err := h.svc.UpdateVenue(c.Request().Context(), &venuepb.UpdateVenueRequest{
		Id: c.Param("id"), Name: req.Name, Address: req.Address, Phone: req.Phone, Email: req.Email,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *VenueHandler) DeleteVenue(c echo.Context) error {
	venueID := c.Param("id")
	log.Info().Str("venue_id", venueID).Msg("Deleting venue")
	if err := h.svc.DeleteVenue(c.Request().Context(), venueID); err != nil {
		log.Error().Err(err).Str("venue_id", venueID).Msg("Failed to delete venue")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	log.Info().Str("venue_id", venueID).Msg("Venue deleted successfully")
	return c.NoContent(http.StatusNoContent)
}

func (h *VenueHandler) ListRooms(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 50
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	resp, err := h.svc.ListRooms(c.Request().Context(), c.Param("venueId"), int32(limit), int32(offset))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *VenueHandler) GetRoom(c echo.Context) error {
	resp, err := h.svc.GetRoom(c.Request().Context(), c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *VenueHandler) CreateRoom(c echo.Context) error {
	var req struct{ Name string }
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	resp, err := h.svc.CreateRoom(c.Request().Context(), &venuepb.CreateRoomRequest{
		VenueId: c.Param("venueId"), Name: req.Name,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, resp)
}

func (h *VenueHandler) UpdateRoom(c echo.Context) error {
	var req struct{ Name string }
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	resp, err := h.svc.UpdateRoom(c.Request().Context(), &venuepb.UpdateRoomRequest{
		Id: c.Param("id"), Name: req.Name,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *VenueHandler) DeleteRoom(c echo.Context) error {
	if err := h.svc.DeleteRoom(c.Request().Context(), c.Param("id")); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *VenueHandler) ListTables(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 50
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	resp, err := h.svc.ListTables(c.Request().Context(), c.Param("roomId"), int32(limit), int32(offset))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *VenueHandler) GetTable(c echo.Context) error {
	resp, err := h.svc.GetTable(c.Request().Context(), c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *VenueHandler) CreateTable(c echo.Context) error {
	var req struct {
		Name     string
		Capacity int32
		CanMerge bool
		Zone     string
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	resp, err := h.svc.CreateTable(c.Request().Context(), &venuepb.CreateTableRequest{
		RoomId: c.Param("roomId"), Name: req.Name, Capacity: req.Capacity, CanMerge: req.CanMerge, Zone: req.Zone,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, resp)
}

func (h *VenueHandler) UpdateTable(c echo.Context) error {
	var req struct {
		Name     string
		Capacity int32
		CanMerge bool
		Zone     string
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	resp, err := h.svc.UpdateTable(c.Request().Context(), &venuepb.UpdateTableRequest{
		Id: c.Param("id"), Name: req.Name, Capacity: req.Capacity, CanMerge: req.CanMerge, Zone: req.Zone,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *VenueHandler) DeleteTable(c echo.Context) error {
	if err := h.svc.DeleteTable(c.Request().Context(), c.Param("id")); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *VenueHandler) GetOpeningHours(c echo.Context) error {
	resp, err := h.svc.GetOpeningHours(c.Request().Context(), c.Param("venueId"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *VenueHandler) SetOpeningHours(c echo.Context) error {
	var req struct {
		Days []struct {
			Weekday   int32
			OpenTime  string
			CloseTime string
		}
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	days := make([]*venuepb.DayHours, len(req.Days))
	for i, d := range req.Days {
		days[i] = &venuepb.DayHours{Weekday: d.Weekday, OpenTime: d.OpenTime, CloseTime: d.CloseTime}
	}
	resp, err := h.svc.SetOpeningHours(c.Request().Context(), &venuepb.SetOpeningHoursRequest{
		VenueId: c.Param("venueId"), Days: days,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *VenueHandler) SetSpecialHours(c echo.Context) error {
	var req struct {
		Date, OpenTime, CloseTime string
		IsClosed                  bool
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	resp, err := h.svc.SetSpecialHours(c.Request().Context(), &venuepb.SetSpecialHoursRequest{
		VenueId: c.Param("venueId"), Date: req.Date, OpenTime: req.OpenTime, CloseTime: req.CloseTime, IsClosed: req.IsClosed,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *VenueHandler) CheckAvailability(c echo.Context) error {
	var req struct {
		VenueID   string `json:"venue_id"`
		Slot      struct {
			Date            string `json:"date"`
			StartTime       string `json:"start_time"`
			DurationMinutes int32  `json:"duration_minutes"`
		} `json:"slot"`
		PartySize int32 `json:"party_size"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	resp, err := h.svc.CheckAvailability(c.Request().Context(), &venuepb.CheckAvailabilityRequest{
		VenueId: req.VenueID,
		Slot: &commonpb.Slot{Date: req.Slot.Date, StartTime: req.Slot.StartTime, DurationMinutes: req.Slot.DurationMinutes},
		PartySize: req.PartySize,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	marshaler := protojson.MarshalOptions{EmitUnpopulated: true, UseProtoNames: true}
	jsonBytes, err := marshaler.Marshal(resp)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.Blob(http.StatusOK, echo.MIMEApplicationJSON, jsonBytes)
}

