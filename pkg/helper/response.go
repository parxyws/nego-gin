package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/shared/domain"
	"github.com/parxyws/nego-gin/pkg/util"
	"github.com/parxyws/nego-gin/pkg/validator"
)

type ApiResponse[T any] struct {
	Success  bool             `json:"success"`
	Message  string           `json:"message"`
	Data     T                `json:"data,omitempty"`
	Error    any              `json:"error,omitempty"`
	Metadata *domain.Metadata `json:"metadata,omitempty"`
}

func Success(c *gin.Context, status int, message string, data any) {
	c.JSON(status, ApiResponse[any]{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, status int, message string, err error) {
	finalStatus := status
	var errDetail any = nil
	var detailData any = nil

	if err != nil {
		errDetail = err.Error()

		translatedErrs := validator.TranslateValidationError(err)
		if len(translatedErrs) > 0 {
			detailData = translatedErrs
		}

		if appErr, ok := util.IsAppError(err); ok {
			if status == http.StatusInternalServerError || status == 0 {
				finalStatus = appErr.Code
			}
			message = appErr.Message
		}
	}

	c.JSON(finalStatus, ApiResponse[any]{
		Success: false,
		Message: message,
		Error:   errDetail,
		Data:    detailData,
	})
}
