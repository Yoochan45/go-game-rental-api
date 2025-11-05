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

type PartnerHandler struct {
	partnerService service.PartnerApplicationService
	bookingService service.BookingService
	validate       *validator.Validate
}

func NewPartnerHandler(partnerService service.PartnerApplicationService, bookingService service.BookingService) *PartnerHandler {
	return &PartnerHandler{
		partnerService: partnerService,
		bookingService: bookingService,
		validate:       validator.New(),
	}
}

func (h *PartnerHandler) ApplyPartner(c echo.Context) error {
	userID := echomw.CurrentUserID(c)
	if userID == 0 {
		return myResponse.Unauthorized(c, "Unauthorized")
	}

	var req dto.CreatePartnerApplicationRequest
	if err := c.Bind(&req); err != nil {
		return myResponse.BadRequest(c, "Invalid input: "+err.Error())
	}
	if err := h.validate.Struct(&req); err != nil {
		return myResponse.BadRequest(c, "Validation error: "+err.Error())
	}

	applicationData := &model.PartnerApplication{
		BusinessName:        req.BusinessName,
		BusinessAddress:     req.BusinessAddress,
		BusinessPhone:       &req.BusinessPhone,
		BusinessDescription: &req.BusinessDescription,
	}

	err := h.partnerService.SubmitApplication(userID, applicationData)
	if err != nil {
		return myResponse.BadRequest(c, err.Error())
	}

	return myResponse.Created(c, "Partner application submitted successfully", nil)
}

func (h *PartnerHandler) GetPartnerBookings(c echo.Context) error {
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

	bookings, err := h.bookingService.GetPartnerBookings(userID, limit, (page-1)*limit)
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

	return myResponse.Paginated(c, "Partner bookings retrieved successfully", bookingDTOs, meta)
}

func (h *PartnerHandler) ConfirmHandover(c echo.Context) error {
	userID := echomw.CurrentUserID(c)
	if userID == 0 {
		return myResponse.Unauthorized(c, "Unauthorized")
	}

	bookingID := myRequest.PathParamUint(c, "id")
	if bookingID == 0 {
		return myResponse.BadRequest(c, "Invalid booking ID")
	}

	err := h.bookingService.ConfirmHandover(userID, bookingID)
	if err != nil {
		return myResponse.BadRequest(c, err.Error())
	}

	return myResponse.Success(c, "Handover confirmed successfully", nil)
}

func (h *PartnerHandler) ConfirmReturn(c echo.Context) error {
	userID := echomw.CurrentUserID(c)
	if userID == 0 {
		return myResponse.Unauthorized(c, "Unauthorized")
	}

	bookingID := myRequest.PathParamUint(c, "id")
	if bookingID == 0 {
		return myResponse.BadRequest(c, "Invalid booking ID")
	}

	err := h.bookingService.ConfirmReturn(userID, bookingID)
	if err != nil {
		return myResponse.BadRequest(c, err.Error())
	}

	return myResponse.Success(c, "Return confirmed successfully", nil)
}
