package handler

import (
	echomw "github.com/Yoochan45/go-api-utils/pkg-echo/middleware"
	myRequest "github.com/Yoochan45/go-api-utils/pkg-echo/request"
	myResponse "github.com/Yoochan45/go-api-utils/pkg-echo/response"
	"github.com/Yoochan45/go-game-rental-api/internal/model"
	"github.com/Yoochan45/go-game-rental-api/internal/model/dto"
	"github.com/Yoochan45/go-game-rental-api/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type BookingHandler struct {
	bookingService service.BookingService
	validate       *validator.Validate
}

func NewBookingHandler(bookingService service.BookingService) *BookingHandler {
	return &BookingHandler{
		bookingService: bookingService,
		validate:       validator.New(),
	}
}

// User endpoints
func (h *BookingHandler) CreateBooking(c echo.Context) error {
	userID := echomw.CurrentUserID(c)
	if userID == 0 {
		return myResponse.Unauthorized(c, "Unauthorized")
	}

	var req dto.CreateBookingRequest
	if err := c.Bind(&req); err != nil {
		return myResponse.BadRequest(c, "Invalid input: "+err.Error())
	}
	if err := h.validate.Struct(&req); err != nil {
		return myResponse.BadRequest(c, "Validation error: "+err.Error())
	}

	bookingData := &model.Booking{
		UserID:    userID,
		GameID:    req.GameID,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Notes:     &req.Notes,
	}

	err := h.bookingService.CreateBooking(userID, bookingData)
	if err != nil {
		return myResponse.BadRequest(c, err.Error())
	}

	return myResponse.Created(c, "Booking created successfully", nil)
}

func (h *BookingHandler) GetMyBookings(c echo.Context) error {
	userID := echomw.CurrentUserID(c)
	if userID == 0 {
		return myResponse.Unauthorized(c, "Unauthorized")
	}

	page := myRequest.QueryInt(c, "page", 1)
	limit := myRequest.QueryInt(c, "limit", 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	bookings, err := h.bookingService.GetUserBookings(userID, limit, (page-1)*limit)
	if err != nil {
		return myResponse.InternalServerError(c, "Failed to retrieve bookings")
	}

	bookingDTOs := dto.ToBookingDTOList(bookings)

	totalCount := int64(len(bookings))
	meta := map[string]any{
		"page":        page,
		"limit":       limit,
		"total":       totalCount,
		"total_pages": (totalCount + int64(limit) - 1) / int64(limit),
	}

	return myResponse.Paginated(c, "Bookings retrieved successfully", bookingDTOs, meta)
}

func (h *BookingHandler) GetBookingDetail(c echo.Context) error {
	userID := echomw.CurrentUserID(c)
	if userID == 0 {
		return myResponse.Unauthorized(c, "Unauthorized")
	}

	bookingID := myRequest.PathParamUint(c, "id")
	if bookingID == 0 {
		return myResponse.BadRequest(c, "Invalid booking ID")
	}

	booking, err := h.bookingService.GetBookingDetail(userID, bookingID)
	if err != nil {
		return myResponse.NotFound(c, err.Error())
	}

	response := dto.ToBookingDTO(booking)
	return myResponse.Success(c, "Booking retrieved successfully", response)
}

func (h *BookingHandler) CancelBooking(c echo.Context) error {
	userID := echomw.CurrentUserID(c)
	if userID == 0 {
		return myResponse.Unauthorized(c, "Unauthorized")
	}

	bookingID := myRequest.PathParamUint(c, "id")
	if bookingID == 0 {
		return myResponse.BadRequest(c, "Invalid booking ID")
	}

	err := h.bookingService.CancelBooking(userID, bookingID)
	if err != nil {
		return myResponse.BadRequest(c, err.Error())
	}

	return myResponse.Success(c, "Booking cancelled successfully", nil)
}

// Admin endpoints
func (h *BookingHandler) GetAllBookings(c echo.Context) error {
	page := myRequest.QueryInt(c, "page", 1)
	limit := myRequest.QueryInt(c, "limit", 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	role := echomw.CurrentRole(c)

	bookings, err := h.bookingService.GetAllBookings(model.UserRole(role), limit, (page-1)*limit)
	if err != nil {
		return myResponse.Forbidden(c, err.Error())
	}

	bookingDTOs := dto.ToBookingDTOList(bookings)

	totalCount := int64(len(bookings))
	meta := map[string]any{
		"page":        page,
		"limit":       limit,
		"total":       totalCount,
		"total_pages": (totalCount + int64(limit) - 1) / int64(limit),
	}

	return myResponse.Paginated(c, "Bookings retrieved successfully", bookingDTOs, meta)
}

func (h *BookingHandler) GetBookingDetailAdmin(c echo.Context) error {
	bookingID := myRequest.PathParamUint(c, "id")
	if bookingID == 0 {
		return myResponse.BadRequest(c, "Invalid booking ID")
	}

	booking, err := h.bookingService.GetBookingDetail(0, bookingID)
	if err != nil {
		return myResponse.NotFound(c, err.Error())
	}

	response := dto.ToBookingDTO(booking)
	return myResponse.Success(c, "Booking retrieved successfully", response)
}
