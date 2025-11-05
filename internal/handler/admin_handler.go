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

type AdminHandler struct {
	partnerService service.PartnerApplicationService
	gameService    service.GameService
	validate       *validator.Validate
}

func NewAdminHandler(partnerService service.PartnerApplicationService, gameService service.GameService) *AdminHandler {
	return &AdminHandler{
		partnerService: partnerService,
		gameService:    gameService,
		validate:       validator.New(),
	}
}

func (h *AdminHandler) GetPartnerApplications(c echo.Context) error {
	page := myRequest.QueryInt(c, "page", 1)
	limit := myRequest.QueryInt(c, "limit", 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	role := echomw.CurrentRole(c)
	applications, err := h.partnerService.GetAllApplications(model.UserRole(role), limit, (page-1)*limit)
	if err != nil {
		return myResponse.Forbidden(c, err.Error())
	}

	applicationDTOs := dto.ToPartnerApplicationDTOList(applications)

	totalCount := int64(len(applications))
	meta := map[string]any{
		"page":        page,
		"limit":       limit,
		"total":       totalCount,
		"total_pages": (totalCount + int64(limit) - 1) / int64(limit),
	}

	return myResponse.Paginated(c, "Partner applications retrieved successfully", applicationDTOs, meta)
}

func (h *AdminHandler) ApprovePartnerApplication(c echo.Context) error {
	applicationID := myRequest.PathParamUint(c, "id")
	if applicationID == 0 {
		return myResponse.BadRequest(c, "Invalid application ID")
	}

	adminID := echomw.CurrentUserID(c)
	role := echomw.CurrentRole(c)

	err := h.partnerService.ApproveApplication(adminID, model.UserRole(role), applicationID)
	if err != nil {
		return myResponse.BadRequest(c, err.Error())
	}

	return myResponse.Success(c, "Partner application approved successfully", nil)
}

func (h *AdminHandler) RejectPartnerApplication(c echo.Context) error {
	applicationID := myRequest.PathParamUint(c, "id")
	if applicationID == 0 {
		return myResponse.BadRequest(c, "Invalid application ID")
	}

	var req struct {
		RejectionReason string `json:"rejection_reason" validate:"required"`
	}
	if err := c.Bind(&req); err != nil {
		return myResponse.BadRequest(c, "Invalid input: "+err.Error())
	}
	if err := h.validate.Struct(&req); err != nil {
		return myResponse.BadRequest(c, "Validation error: "+err.Error())
	}

	adminID := echomw.CurrentUserID(c)
	role := echomw.CurrentRole(c)

	err := h.partnerService.RejectApplication(adminID, model.UserRole(role), applicationID, req.RejectionReason)
	if err != nil {
		return myResponse.BadRequest(c, err.Error())
	}

	return myResponse.Success(c, "Partner application rejected successfully", nil)
}

func (h *AdminHandler) GetGameListings(c echo.Context) error {
	page := myRequest.QueryInt(c, "page", 1)
	limit := myRequest.QueryInt(c, "limit", 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	role := echomw.CurrentRole(c)
	games, err := h.gameService.GetAllGames(model.UserRole(role), limit, (page-1)*limit)
	if err != nil {
		return myResponse.Forbidden(c, err.Error())
	}

	gameDTOs := dto.ToGameDTOList(games)

	totalCount := int64(len(games))
	meta := map[string]any{
		"page":        page,
		"limit":       limit,
		"total":       totalCount,
		"total_pages": (totalCount + int64(limit) - 1) / int64(limit),
	}

	return myResponse.Paginated(c, "Game listings retrieved successfully", gameDTOs, meta)
}

func (h *AdminHandler) ApproveGameListing(c echo.Context) error {
	gameID := myRequest.PathParamUint(c, "id")
	if gameID == 0 {
		return myResponse.BadRequest(c, "Invalid game ID")
	}

	adminID := echomw.CurrentUserID(c)
	role := echomw.CurrentRole(c)

	err := h.gameService.ApproveGame(adminID, model.UserRole(role), gameID)
	if err != nil {
		return myResponse.BadRequest(c, err.Error())
	}

	return myResponse.Success(c, "Game listing approved successfully", nil)
}
