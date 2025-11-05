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

// GetAllGames godoc
// @Summary Get all games
// @Description Get list of all available games (public)
// @Tags Games
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} map[string]interface{} "Games retrieved successfully"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /games [get]
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

// GetGameDetail godoc
// @Summary Get game detail
// @Description Get detailed information about a specific game
// @Tags Games
// @Accept json
// @Produce json
// @Param id path int true "Game ID"
// @Success 200 {object} dto.GameDTO "Game retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid game ID"
// @Failure 404 {object} map[string]interface{} "Game not found"
// @Router /games/{id} [get]
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

// SearchGames godoc
// @Summary Search games
// @Description Search games by name, description, or platform
// @Tags Games
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} map[string]interface{} "Games search results"
// @Failure 400 {object} map[string]interface{} "Search query required"
// @Router /games/search [get]
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

// CreateGame godoc
// @Summary Create new game
// @Description Create a new game listing (Admin only)
// @Tags Admin - Games
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateGameRequest true "Game details"
// @Success 201 {object} map[string]interface{} "Game created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Router /admin/games [post]
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

// UpdateGame godoc
// @Summary Update game
// @Description Update game information (Admin only)
// @Tags Admin - Games
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Game ID"
// @Param request body dto.UpdateGameRequest true "Updated game details"
// @Success 200 {object} map[string]interface{} "Game updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Router /admin/games/{id} [put]
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

// DeleteGame godoc
// @Summary Delete game
// @Description Delete a game listing (Admin only)
// @Tags Admin - Games
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Game ID"
// @Success 200 {object} map[string]interface{} "Game deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid game ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Router /admin/games/{id} [delete]
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
