package utils

import (
	"strings"

	myResponse "github.com/Yoochan45/go-api-utils/pkg-echo/response"
	"github.com/labstack/echo/v4"
)

// MapServiceError maps service errors to appropriate HTTP responses
func MapServiceError(c echo.Context, err error) error {
	errMsg := err.Error()
	
	switch {
	case strings.Contains(errMsg, "not found"):
		return myResponse.NotFound(c, errMsg)
	
	case strings.Contains(errMsg, "not owned") || 
		 strings.Contains(errMsg, "insufficient permission"):
		return myResponse.Forbidden(c, errMsg)
	
	case strings.Contains(errMsg, "cannot cancel") || 
		 strings.Contains(errMsg, "cannot confirm"):
		return myResponse.BadRequest(c, errMsg)
	
	default:
		return myResponse.BadRequest(c, errMsg)
	}
}