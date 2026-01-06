package route

import (
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/user"
)

func UserRoute(route *gin.Engine, controller user.UserController) {
	user := route.Group("/user")
	{
		user.POST("")
	}
}
