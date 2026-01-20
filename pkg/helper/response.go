package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/shared/domain"
	"github.com/parxyws/nego-gin/pkg/util"
	"github.com/parxyws/nego-gin/pkg/validator"
	"github.com/sirupsen/logrus"
)

type ApiResponse[T any] struct {
	Success  bool             `json:"success"`
	Message  string           `json:"message"`
	Data     T                `json:"data,omitempty"`
	Details  any              `json:"details,omitempty"`
	Metadata *domain.Metadata `json:"metadata,omitempty"`
}

func Success(ctx *gin.Context, status int, message string, data any) {
	logrus.WithFields(logrus.Fields{
		"status": status,
		"path":   ctx.Request.URL.Path,
		"method": ctx.Request.Method,
	}).Info(message)

	ctx.JSON(status, ApiResponse[any]{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(ctx *gin.Context, status int, message string, err error) {
	statusCode := status
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
				statusCode = appErr.Code
			}
			message = appErr.Message
		}
	}

	logFields := logrus.Fields{
		"status": statusCode,
		"path":   ctx.Request.URL.Path,
		"method": ctx.Request.Method,
		"error":  errDetail,
	}

	if statusCode >= 500 {
		logrus.WithFields(logFields).Error(message)
	} else {
		logrus.WithFields(logFields).Warn(message)
	}

	ctx.JSON(statusCode, ApiResponse[any]{
		Success: false,
		Message: message,
		Details: errDetail,
		Data:    detailData,
	})
}
