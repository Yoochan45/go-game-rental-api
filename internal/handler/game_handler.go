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
	"github.com/lib/pq"
)

type GameHandler struct {
	gameService service.GameService
	validate    *validator.Validate
}

func NewGameHandler(gameService service.GameService) *GameHandler {
	return &GameHandler{
		gameService: gameService,
		validate:    validator.New(),
	}
}

// Public endpoints
func (h *GameHandler) GetAllGames(c echo.Context) error {
	page := myRequest.QueryInt(c, "page", 1)
	limit := myRequest.QueryInt(c, "limit", 10)

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Use "customer" role for public access
	games, err := h.gameService.GetAllGames("customer", limit, (page-1)*limit)
	if err != nil {
		return myResponse.InternalServerError(c, "Failed to retrieve games")
	}

	gameDTOs := dto.ToGameDTOList(games)

	totalCount := int64(len(games))
	meta := map[string]any{
		"page":        page,
		"limit":       limit,
		"total":       totalCount,
		"total_pages": (totalCount + int64(limit) - 1) / int64(limit),
	}

	return myResponse.Paginated(c, "Games retrieved successfully", gameDTOs, meta)
}

func (h *GameHandler) GetGameDetail(c echo.Context) error {
	gameID := myRequest.PathParamUint(c, "id")
	if gameID == 0 {
		return myResponse.BadRequest(c, "Invalid game ID")
	}

	game, err := h.gameService.GetGameDetail(gameID)
	if err != nil {
		return myResponse.NotFound(c, "Game not found")
	}

	response := dto.ToGameDTO(game)
	return myResponse.Success(c, "Game retrieved successfully", response)
}

func (h *GameHandler) SearchGames(c echo.Context) error {
	query := myRequest.QueryString(c, "q", "")
	if query == "" {
		return myResponse.BadRequest(c, "Search query is required")
	}

	page := myRequest.QueryInt(c, "page", 1)
	limit := myRequest.QueryInt(c, "limit", 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	games, err := h.gameService.SearchGames(query, limit, (page-1)*limit)
	if err != nil {
		return myResponse.InternalServerError(c, "Failed to search games")
	}

	gameDTOs := dto.ToGameDTOList(games)

	totalCount := int64(len(games))
	meta := map[string]any{
		"page":        page,
		"limit":       limit,
		"total":       totalCount,
		"total_pages": (totalCount + int64(limit) - 1) / int64(limit),
	}

	return myResponse.Paginated(c, "Games search results", gameDTOs, meta)
}

// Admin endpoints
func (h *GameHandler) CreateGame(c echo.Context) error {
	var req dto.CreateGameRequest
	if err := c.Bind(&req); err != nil {
		return myResponse.BadRequest(c, "Invalid input: "+err.Error())
	}
	if err := h.validate.Struct(&req); err != nil {
		return myResponse.BadRequest(c, "Validation error: "+err.Error())
	}

	adminID := echomw.CurrentUserID(c)

	gameData := &model.Game{
		CategoryID:        req.CategoryID,
		Name:              req.Name,
		Description:       &req.Description,
		Platform:          &req.Platform,
		Stock:             req.Stock,
		RentalPricePerDay: req.RentalPricePerDay, // Sesuai dengan field di DTO
		SecurityDeposit:   req.SecurityDeposit,
		Condition:         req.Condition,
		Images:            pq.StringArray(req.Images), // Sesuai dengan field di DTO
	}

	// CreateGame takes (adminID, gameData)
	err := h.gameService.CreateGame(adminID, gameData)
	if err != nil {
		return myResponse.BadRequest(c, err.Error())
	}

	return myResponse.Created(c, "Game created successfully", nil)
}

func (h *GameHandler) UpdateGame(c echo.Context) error {
	gameID := myRequest.PathParamUint(c, "id")
	if gameID == 0 {
		return myResponse.BadRequest(c, "Invalid game ID")
	}

	var req dto.UpdateGameRequest
	if err := c.Bind(&req); err != nil {
		return myResponse.BadRequest(c, "Invalid input: "+err.Error())
	}
	if err := h.validate.Struct(&req); err != nil {
		return myResponse.BadRequest(c, "Validation error: "+err.Error())
	}

	adminID := echomw.CurrentUserID(c)

	updateData := &model.Game{
		CategoryID:        req.CategoryID,
		Name:              req.Name,
		Description:       &req.Description,
		Platform:          &req.Platform,
		RentalPricePerDay: req.RentalPricePerDay,
		SecurityDeposit:   req.SecurityDeposit,
		Condition:         req.Condition,
		Images:            pq.StringArray(req.Images),
	}

	err := h.gameService.UpdateGame(adminID, gameID, updateData)
	if err != nil {
		return myResponse.BadRequest(c, err.Error())
	}

	return myResponse.Success(c, "Game updated successfully", nil)
}

func (h *GameHandler) DeleteGame(c echo.Context) error {
	gameID := myRequest.PathParamUint(c, "id")
	if gameID == 0 {
		return myResponse.BadRequest(c, "Invalid game ID")
	}

	adminID := echomw.CurrentUserID(c)
	err := h.gameService.DeleteGame(adminID, gameID)
	if err != nil {
		return myResponse.BadRequest(c, err.Error())
	}

	return myResponse.Success(c, "Game deleted successfully", nil)
}

func (h *GameHandler) GetAllGamesAdmin(c echo.Context) error {
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

	return myResponse.Paginated(c, "Games retrieved successfully", gameDTOs, meta)
}
