package main

import (
	"log"
	"os"

	myOrm "github.com/Yoochan45/go-api-utils/pkg-echo/orm"
	myConfig "github.com/Yoochan45/go-api-utils/pkg/config"
	"github.com/Yoochan45/go-game-rental-api/app/echo-server/router"
	"github.com/Yoochan45/go-game-rental-api/internal/handler"
	"github.com/Yoochan45/go-game-rental-api/internal/model"
	"github.com/Yoochan45/go-game-rental-api/internal/repository"
	"github.com/Yoochan45/go-game-rental-api/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := myConfig.LoadEnv()
	JwtSecret := os.Getenv("JWT_SECRET")
	if JwtSecret == "" {
		JwtSecret = "dev-secret"
	}

	db, err := myOrm.Init(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate all models
	err = db.AutoMigrate(
		&model.User{},
		&model.RefreshToken{},
		&model.Category{},
		&model.Game{},
		&model.Booking{},
		&model.Payment{},
		&model.Review{},
		&model.PartnerApplication{},
		&model.Dispute{},
	)
	if err != nil {
		log.Println("Migration warning:", err)
		// Continue even if tables already exist
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	gameRepo := repository.NewGameRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	partnerRepo := repository.NewPartnerApplicationRepository(db)
	disputeRepo := repository.NewDisputeRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	gameService := service.NewGameService(gameRepo, userRepo)
	bookingService := service.NewBookingService(bookingRepo, gameRepo, userRepo)
	paymentService := service.NewPaymentService(paymentRepo, bookingRepo, userRepo, bookingService)
	reviewService := service.NewReviewService(reviewRepo, bookingRepo)
	partnerService := service.NewPartnerApplicationService(partnerRepo, userRepo)
	disputeService := service.NewDisputeService(disputeRepo, bookingRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(userService, JwtSecret)
	userHandler := handler.NewUserHandler(userService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	gameHandler := handler.NewGameHandler(gameService)
	bookingHandler := handler.NewBookingHandler(bookingService)
	paymentHandler := handler.NewPaymentHandler(paymentService)
	reviewHandler := handler.NewReviewHandler(reviewService)
	partnerHandler := handler.NewPartnerHandler(partnerService, bookingService)
	adminHandler := handler.NewAdminHandler(partnerService, gameService)
	disputeHandler := handler.NewDisputeHandler(disputeService)

	// Setup Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Register routes
	router.RegisterRoutes(
		e,
		authHandler,
		userHandler,
		categoryHandler,
		gameHandler,
		bookingHandler,
		paymentHandler,
		reviewHandler,
		partnerHandler,
		adminHandler,
		disputeHandler,
		JwtSecret,
	)

	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	log.Fatal(e.Start(":" + port))
}
