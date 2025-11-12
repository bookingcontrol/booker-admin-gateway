package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	venuepb "github.com/bookingcontrol/booker-contracts-go/venue"
)

// Venue handlers
func (h *Handler) ListVenues(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 50
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	resp, err := h.venueUseCase.ListVenues(c.Request().Context(), int32(limit), int32(offset))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetVenue(c echo.Context) error {
	resp, err := h.venueUseCase.GetVenue(c.Request().Context(), c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) CreateVenue(c echo.Context) error {
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

	log.Info().
		Str("name", req.Name).
		Str("timezone", req.Timezone).
		Str("address", req.Address).
		Str("phone", req.Phone).
		Str("email", req.Email).
		Msg("Creating venue")

	resp, err := h.venueUseCase.CreateVenue(c.Request().Context(), &venuepb.CreateVenueRequest{
		Name:     req.Name,
		Timezone: req.Timezone,
		Address:  req.Address,
		Phone:    req.Phone,
		Email:    req.Email,
	})
	if err != nil {
		log.Error().Err(err).
			Str("name", req.Name).
			Msg("Failed to create venue")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	log.Info().
		Str("venue_id", resp.Id).
		Str("name", resp.Name).
		Msg("Venue created successfully")

	return c.JSON(http.StatusCreated, resp)
}

func (h *Handler) UpdateVenue(c echo.Context) error {
	var req struct {
		Name    string `json:"name"`
		Address string `json:"address"`
		Phone   string `json:"phone"`
		Email   string `json:"email"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	resp, err := h.venueUseCase.UpdateVenue(c.Request().Context(), &venuepb.UpdateVenueRequest{
		Id:      c.Param("id"),
		Name:    req.Name,
		Address: req.Address,
		Phone:   req.Phone,
		Email:   req.Email,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) DeleteVenue(c echo.Context) error {
	venueID := c.Param("id")

	log.Info().Str("venue_id", venueID).Msg("Deleting venue")

	if err := h.venueUseCase.DeleteVenue(c.Request().Context(), venueID); err != nil {
		log.Error().Err(err).
			Str("venue_id", venueID).
			Msg("Failed to delete venue")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	log.Info().Str("venue_id", venueID).Msg("Venue deleted successfully")
	return c.NoContent(http.StatusNoContent)
}

// Room handlers
func (h *Handler) ListRooms(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 50
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	resp, err := h.venueUseCase.ListRooms(c.Request().Context(), c.Param("venueId"), int32(limit), int32(offset))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetRoom(c echo.Context) error {
	resp, err := h.venueUseCase.GetRoom(c.Request().Context(), c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) CreateRoom(c echo.Context) error {
	var req struct {
		Name string `json:"name"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	resp, err := h.venueUseCase.CreateRoom(c.Request().Context(), &venuepb.CreateRoomRequest{
		VenueId: c.Param("venueId"),
		Name:    req.Name,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h *Handler) UpdateRoom(c echo.Context) error {
	var req struct {
		Name string `json:"name"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	resp, err := h.venueUseCase.UpdateRoom(c.Request().Context(), &venuepb.UpdateRoomRequest{
		Id:   c.Param("id"),
		Name: req.Name,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) DeleteRoom(c echo.Context) error {
	if err := h.venueUseCase.DeleteRoom(c.Request().Context(), c.Param("id")); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// Table handlers
func (h *Handler) ListTables(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 50
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	resp, err := h.venueUseCase.ListTables(c.Request().Context(), c.Param("roomId"), int32(limit), int32(offset))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetTable(c echo.Context) error {
	resp, err := h.venueUseCase.GetTable(c.Request().Context(), c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) CreateTable(c echo.Context) error {
	var req struct {
		Name     string `json:"name"`
		Capacity int32  `json:"capacity"`
		CanMerge bool   `json:"can_merge"`
		Zone     string `json:"zone"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	resp, err := h.venueUseCase.CreateTable(c.Request().Context(), &venuepb.CreateTableRequest{
		RoomId:   c.Param("roomId"),
		Name:     req.Name,
		Capacity: req.Capacity,
		CanMerge: req.CanMerge,
		Zone:     req.Zone,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h *Handler) UpdateTable(c echo.Context) error {
	var req struct {
		Name     string `json:"name"`
		Capacity int32  `json:"capacity"`
		CanMerge bool   `json:"can_merge"`
		Zone     string `json:"zone"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	resp, err := h.venueUseCase.UpdateTable(c.Request().Context(), &venuepb.UpdateTableRequest{
		Id:       c.Param("id"),
		Name:     req.Name,
		Capacity: req.Capacity,
		CanMerge: req.CanMerge,
		Zone:     req.Zone,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) DeleteTable(c echo.Context) error {
	if err := h.venueUseCase.DeleteTable(c.Request().Context(), c.Param("id")); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// Schedule handlers
func (h *Handler) GetOpeningHours(c echo.Context) error {
	resp, err := h.venueUseCase.GetOpeningHours(c.Request().Context(), c.Param("venueId"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) SetOpeningHours(c echo.Context) error {
	var req struct {
		Days []struct {
			Weekday   int32  `json:"weekday"`
			OpenTime  string `json:"open_time"`
			CloseTime string `json:"close_time"`
		} `json:"days"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	days := make([]*venuepb.DayHours, len(req.Days))
	for i, d := range req.Days {
		days[i] = &venuepb.DayHours{
			Weekday:   d.Weekday,
			OpenTime:  d.OpenTime,
			CloseTime: d.CloseTime,
		}
	}

	resp, err := h.venueUseCase.SetOpeningHours(c.Request().Context(), &venuepb.SetOpeningHoursRequest{
		VenueId: c.Param("venueId"),
		Days:    days,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) SetSpecialHours(c echo.Context) error {
	var req struct {
		Date      string `json:"date"`
		OpenTime  string `json:"open_time"`
		CloseTime string `json:"close_time"`
		IsClosed  bool   `json:"is_closed"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	resp, err := h.venueUseCase.SetSpecialHours(c.Request().Context(), &venuepb.SetSpecialHoursRequest{
		VenueId:   c.Param("venueId"),
		Date:      req.Date,
		OpenTime:  req.OpenTime,
		CloseTime: req.CloseTime,
		IsClosed:  req.IsClosed,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

