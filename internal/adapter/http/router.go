package http

import (
	"github.com/labstack/echo/v4"
	"github.com/bookingcontrol/booker-admin-gateway/internal/adapter/http/middleware"
	ucauth "github.com/bookingcontrol/booker-admin-gateway/internal/usecase/auth"
	ucvenue "github.com/bookingcontrol/booker-admin-gateway/internal/usecase/venue"
	ucbooking "github.com/bookingcontrol/booker-admin-gateway/internal/usecase/booking"
)

func SetupRouter(
	authSvc *ucauth.Service,
	venueSvc *ucvenue.Service,
	bookingSvc *ucbooking.Service,
	mw *middleware.Middleware,
) *echo.Echo {
	e := echo.New()
	e.Use(middleware.MetricsMiddleware("admin-gateway"))
	mw.SetupMiddleware(e)

	authH := NewAuthHandler(authSvc)
	venueH := NewVenueHandler(venueSvc)
	bookingH := NewBookingHandler(bookingSvc)

	e.GET("/metrics", bookingH.Metrics)
	e.GET("/api", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"service": "Admin Gateway", "version": "1.0.0",
			"endpoints": map[string]string{
				"auth": "/api/v1/auth/login", "venues": "/api/v1/venues",
				"bookings": "/api/v1/bookings", "availability": "/api/v1/availability/check",
				"websocket": "/api/v1/ws",
			},
		})
	})

	api := e.Group("/api/v1")
	api.POST("/auth/register", authH.Register)
	api.POST("/auth/login", authH.Login)
	api.POST("/auth/refresh", authH.RefreshToken)

	protected := api.Group("", mw.AuthMiddleware)
	protected.GET("/venues", venueH.ListVenues)
	protected.GET("/venues/:id", venueH.GetVenue)
	protected.POST("/venues", venueH.CreateVenue)
	protected.PUT("/venues/:id", venueH.UpdateVenue)
	protected.DELETE("/venues/:id", venueH.DeleteVenue)
	protected.GET("/venues/:venueId/rooms", venueH.ListRooms)
	protected.GET("/rooms/:id", venueH.GetRoom)
	protected.POST("/venues/:venueId/rooms", venueH.CreateRoom)
	protected.PUT("/rooms/:id", venueH.UpdateRoom)
	protected.DELETE("/rooms/:id", venueH.DeleteRoom)
	protected.GET("/rooms/:roomId/tables", venueH.ListTables)
	protected.GET("/tables/:id", venueH.GetTable)
	protected.POST("/rooms/:roomId/tables", venueH.CreateTable)
	protected.PUT("/tables/:id", venueH.UpdateTable)
	protected.DELETE("/tables/:id", venueH.DeleteTable)
	protected.GET("/venues/:venueId/schedule", venueH.GetOpeningHours)
	protected.POST("/venues/:venueId/schedule", venueH.SetOpeningHours)
	protected.POST("/venues/:venueId/special-hours", venueH.SetSpecialHours)
	protected.GET("/bookings", bookingH.ListBookings)
	protected.GET("/bookings/:id", bookingH.GetBooking)
	protected.POST("/bookings", bookingH.CreateBooking)
	protected.POST("/bookings/:id/confirm", bookingH.ConfirmBooking)
	protected.POST("/bookings/:id/cancel", bookingH.CancelBooking)
	protected.POST("/bookings/:id/seat", bookingH.MarkSeated)
	protected.POST("/bookings/:id/finish", bookingH.MarkFinished)
	protected.POST("/bookings/:id/no-show", bookingH.MarkNoShow)
	protected.POST("/availability/check", venueH.CheckAvailability)
	protected.GET("/ws", bookingH.WebSocket)
	e.Static("/", "web/dist")
	return e
}

