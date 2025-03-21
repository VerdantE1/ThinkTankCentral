package middleware

import (
	"ThinkTankCentral/model/appTypes"
	"ThinkTankCentral/model/response"
	"ThinkTankCentral/utils"
	"github.com/gin-gonic/gin"
)

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID := utils.GetRoleID(c)
		if roleID != appTypes.Admin {
			response.Forbidden("Access denied.Admin privileges are required", c)
			c.Abort()
			return
		}
		c.Next()
	}
}
