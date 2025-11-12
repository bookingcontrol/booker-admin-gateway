package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bookingcontrol/booker-admin-gateway/internal/config"
	"github.com/bookingcontrol/booker-admin-gateway/internal/delivery/http/handler"
	"github.com/bookingcontrol/booker-admin-gateway/internal/delivery/http/middleware"
	"github.com/bookingcontrol/booker-admin-gateway/internal/infrastructure/redis"
	"github.com/bookingcontrol/booker-admin-gateway/internal/infrastructure/repository"
	"github.com/bookingcontrol/booker-admin-gateway/internal/infrastructure/tracing"
	"github.com/bookingcontrol/booker-admin-gateway/internal/usecase"
	bookingpb "github.com/bookingcontrol/booker-contracts-go/booking"
	venuepb "github.com/bookingcontrol/booker-contracts-go/venue"
)

func main() {
	cfg := config.Load()

	// Logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if cfg.Env == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Tracing
	shutdown, err := tracing.InitTracer("admin-gateway", cfg.JaegerEndpoint)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize tracer")
	}
	defer shutdown()

	// Redis
	redisClient := redis.NewClient(cfg.RedisAddr, cfg.RedisPassword)

	// gRPC clients
	venueConn, err := grpc.Dial(
		cfg.GRPCVenueAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to venue service")
	}
	defer venueConn.Close()

	bookingConn, err := grpc.Dial(
		cfg.GRPCBookingAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to booking service")
	}
	defer bookingConn.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(redisClient)
	venueRepo := repository.NewVenueRepository(venuepb.NewVenueServiceClient(venueConn))
	bookingRepo := repository.NewBookingRepository(bookingpb.NewBookingServiceClient(bookingConn))

	// Initialize use cases
	authUseCase := usecase.NewAuthUseCase(userRepo)
	venueUseCase := usecase.NewVenueUseCase(venueRepo)
	bookingUseCase := usecase.NewBookingUseCase(bookingRepo)

	// Initialize handlers
	h := handler.New(authUseCase, venueUseCase, bookingUseCase)

	// Initialize middleware
	mw := middleware.New(redisClient, cfg)

	// Setup routes
	e := h.SetupRoutes(mw)

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := e.Start(fmt.Sprintf(":%d", cfg.Port)); err != nil {
			log.Fatal().Err(err).Msg("Server failed")
		}
	}()

	log.Info().Int("port", cfg.Port).Msg("Admin gateway started")

	<-ctx.Done()
	log.Info().Msg("Shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("Server shutdown error")
	}
}
