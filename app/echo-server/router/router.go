package router

import (
	myMiddleware "github.com/Yoochan45/go-api-utils/pkg-echo/middleware"
	"github.com/Yoochan45/go-game-rental-api/internal/handler"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(
	e *echo.Echo,
	authH *handler.AuthHandler,
	userH *handler.UserHandler,
	categoryH *handler.CategoryHandler,
	gameH *handler.GameHandler,
	bookingH *handler.BookingHandler,
	paymentH *handler.PaymentHandler,
	reviewH *handler.ReviewHandler,
	jwtSecret string,
) {
	// ============= Public Endpoints (No Auth) =============
	// Auth
	e.POST("/auth/register", authH.Register)
	e.POST("/auth/login", authH.Login)

	// Public game catalog
	e.GET("/games", gameH.GetAllGames)
	e.GET("/games/:id", gameH.GetGameDetail)
	e.GET("/games/search", gameH.SearchGames)

	// Public categories
	e.GET("/categories", categoryH.GetAllCategories)
	e.GET("/categories/:id", categoryH.GetCategoryDetail)

	// Public game reviews
	e.GET("/games/:game_id/reviews", reviewH.GetGameReviews)

	// Payment webhook (public but validated by provider)
	e.POST("/webhooks/payments", paymentH.PaymentWebhook)

	// ============= Protected Routes (Authenticated Users) =============
	jwtConfig := myMiddleware.JWTConfig{
		SecretKey:      jwtSecret,
		UseCustomToken: false,
	}

	protected := e.Group("")
	protected.Use(myMiddleware.JWTMiddleware(jwtConfig))

	// User profile
	protected.GET("/users/me", userH.GetMyProfile)
	protected.PUT("/users/me", userH.UpdateMyProfile)

	// Customer bookings
	protected.POST("/bookings", bookingH.CreateBooking)
	protected.GET("/bookings/my", bookingH.GetMyBookings)
	protected.GET("/bookings/:booking_id", bookingH.GetBookingDetail)
	protected.PATCH("/bookings/:booking_id/cancel", bookingH.CancelBooking)

	// Payments
	protected.POST("/bookings/:booking_id/payments", paymentH.CreatePayment)
	protected.GET("/bookings/:booking_id/payments", paymentH.GetPaymentByBooking)

	// Reviews
	protected.POST("/bookings/:booking_id/reviews", reviewH.CreateReview)

	// ============= Admin Routes (Admin & Super Admin) =============
	admin := protected.Group("/admin")
	admin.Use(myMiddleware.RequireRoles("admin", "super_admin"))

	// Game management (admin owns games directly)
	admin.POST("/games", gameH.CreateGame)
	admin.PUT("/games/:id", gameH.UpdateGame)
	admin.DELETE("/games/:id", gameH.DeleteGame)

	// Category management
	admin.POST("/categories", categoryH.CreateCategory)
	admin.PUT("/categories/:id", categoryH.UpdateCategory)
	admin.DELETE("/categories/:id", categoryH.DeleteCategory)

	// Booking management
	admin.GET("/bookings", bookingH.GetAllBookings)
	admin.PATCH("/bookings/:id/status", bookingH.UpdateBookingStatus)

	// Payment management
	admin.GET("/payments", paymentH.GetAllPayments)
	admin.GET("/payments/:id", paymentH.GetPaymentDetail)
	admin.GET("/payments/status", paymentH.GetPaymentsByStatus)

	// User management
	admin.GET("/users", userH.GetAllUsers)
	admin.GET("/users/:id", userH.GetUserDetail)
	admin.PATCH("/users/:id/role", userH.UpdateUserRole)
	admin.PATCH("/users/:id/status", userH.ToggleUserStatus)
	admin.DELETE("/users/:id", userH.DeleteUser)
}
