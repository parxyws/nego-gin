package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/shared/dto"
	"github.com/parxyws/nego-gin/pkg/util"
	"github.com/parxyws/nego-gin/pkg/validator"
)

// Success sends a success response with standardized structure
func Success(c *gin.Context, status int, message string, data any) {
	c.JSON(status, dto.ApiResponse[any]{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error sends an error response with standardized structure
// It automatically handles AppError types to set the correct HTTP status code
// and translates validation errors into a map of field messages.
func Error(c *gin.Context, status int, message string, err error) {
	finalStatus := status
	var errDetail any = nil
	var detailData any = nil

	if err != nil {
		errDetail = err.Error()

		// If it's a validation error, translate it
		translatedErrs := validator.TranslateValidationError(err)
		if len(translatedErrs) > 0 {
			detailData = translatedErrs
		}

		// If it's an AppError, use its internal code if status is default (e.g., 500)
		if appErr, ok := util.IsAppError(err); ok {
			if status == http.StatusInternalServerError || status == 0 {
				finalStatus = appErr.Code
			}
			message = appErr.Message
		}
	}

	c.JSON(finalStatus, dto.ApiResponse[any]{
		Success: false,
		Message: message,
		Error:   errDetail,
		Data:    detailData,
	})
}
