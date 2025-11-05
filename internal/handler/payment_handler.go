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

type PaymentHandler struct {
	paymentService service.PaymentService
	validate       *validator.Validate
}

func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		validate:       validator.New(),
	}
}

// User endpoints
func (h *PaymentHandler) CreatePayment(c echo.Context) error {
	userID := echomw.CurrentUserID(c)
	if userID == 0 {
		return myResponse.Unauthorized(c, "Unauthorized")
	}

	bookingID := myRequest.PathParamUint(c, "booking_id")
	if bookingID == 0 {
		return myResponse.BadRequest(c, "Invalid booking ID")
	}

	var req dto.CreatePaymentRequest
	if err := c.Bind(&req); err != nil {
		return myResponse.BadRequest(c, "Invalid input: "+err.Error())
	}
	if err := h.validate.Struct(&req); err != nil {
		return myResponse.BadRequest(c, "Validation error: "+err.Error())
	}

	payment, err := h.paymentService.CreatePayment(userID, bookingID, req.Provider)
	if err != nil {
		return myResponse.BadRequest(c, err.Error())
	}

	response := dto.ToPaymentDTO(payment)
	return myResponse.Created(c, "Payment created successfully", response)
}

func (h *PaymentHandler) GetPaymentByBooking(c echo.Context) error {
	userID := echomw.CurrentUserID(c)
	if userID == 0 {
		return myResponse.Unauthorized(c, "Unauthorized")
	}

	bookingID := myRequest.PathParamUint(c, "booking_id")
	if bookingID == 0 {
		return myResponse.BadRequest(c, "Invalid booking ID")
	}

	payment, err := h.paymentService.GetPaymentByBooking(userID, bookingID)
	if err != nil {
		return myResponse.NotFound(c, err.Error())
	}

	response := dto.ToPaymentDTO(payment)
	return myResponse.Success(c, "Payment retrieved successfully", response)
}

func (h *PaymentHandler) GetPaymentDetail(c echo.Context) error {
	paymentID := myRequest.PathParamUint(c, "id")
	if paymentID == 0 {
		return myResponse.BadRequest(c, "Invalid payment ID")
	}

	role := echomw.CurrentRole(c)
	payment, err := h.paymentService.GetPaymentDetail(model.UserRole(role), paymentID)
	if err != nil {
		return myResponse.NotFound(c, err.Error())
	}

	response := dto.ToPaymentDTO(payment)
	return myResponse.Success(c, "Payment retrieved successfully", response)
}

// Webhook endpoint (public)
func (h *PaymentHandler) PaymentWebhook(c echo.Context) error {
	var req dto.PaymentWebhookRequest
	if err := c.Bind(&req); err != nil {
		return myResponse.BadRequest(c, "Invalid webhook payload: "+err.Error())
	}
	if err := h.validate.Struct(&req); err != nil {
		return myResponse.BadRequest(c, "Webhook validation error: "+err.Error())
	}

	// Handle pointer fields properly
	var paymentMethod string
	if req.PaymentMethod != nil {
		paymentMethod = *req.PaymentMethod
	}

	// ProcessWebhook takes individual parameters based on actual interface
	err := h.paymentService.ProcessWebhook(req.ProviderPaymentID, req.Status, paymentMethod, req.FailureReason)
	if err != nil {
		return myResponse.BadRequest(c, err.Error())
	}

	return myResponse.Success(c, "Webhook processed successfully", nil)
}

// Admin endpoints
func (h *PaymentHandler) GetAllPayments(c echo.Context) error {
	page := myRequest.QueryInt(c, "page", 1)
	limit := myRequest.QueryInt(c, "limit", 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	role := echomw.CurrentRole(c)
	payments, err := h.paymentService.GetAllPayments(model.UserRole(role), limit, (page-1)*limit)
	if err != nil {
		return myResponse.Forbidden(c, err.Error())
	}

	paymentDTOs := dto.ToPaymentDTOList(payments)

	totalCount := int64(len(payments))
	meta := map[string]any{
		"page":        page,
		"limit":       limit,
		"total":       totalCount,
		"total_pages": (totalCount + int64(limit) - 1) / int64(limit),
	}

	return myResponse.Paginated(c, "Payments retrieved successfully", paymentDTOs, meta)
}

func (h *PaymentHandler) GetPaymentsByStatus(c echo.Context) error {
	page := myRequest.QueryInt(c, "page", 1)
	limit := myRequest.QueryInt(c, "limit", 10)
	status := c.QueryParam("status")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	role := echomw.CurrentRole(c)
	payments, err := h.paymentService.GetPaymentsByStatus(model.UserRole(role), model.PaymentStatus(status), limit, (page-1)*limit)
	if err != nil {
		return myResponse.Forbidden(c, err.Error())
	}

	paymentDTOs := dto.ToPaymentDTOList(payments)

	totalCount := int64(len(payments))
	meta := map[string]any{
		"page":        page,
		"limit":       limit,
		"total":       totalCount,
		"total_pages": (totalCount + int64(limit) - 1) / int64(limit),
	}

	return myResponse.Paginated(c, "Payments retrieved successfully", paymentDTOs, meta)
}
