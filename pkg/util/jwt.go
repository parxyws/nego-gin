package util

import (
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/user/domain"
)

func authorizeHandler(c *gin.Context, data any) bool {
	user, ok := data.(*domain.User)
	if !ok {
		return false
	}

	//path := c.Request.URL.Path
	//method := c.Request.Method

	var _roles []int32
	for i, v := range user.UserRoles {
		_roles[i] = v.RoleID
	}

	if slices.Contains(_roles, 1) {
		return true
	}

	return true
}
