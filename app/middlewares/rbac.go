package middlewares

import (
	"apiok-admin/app/services"
	"apiok-admin/app/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RBAC() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("auth-token")
		path := c.Request.URL.Path
		if services.RBACAllowed(token, path) {
			c.Next()
			return
		}
		utils.CustomError(c, http.StatusForbidden, "无权访问")
		c.Abort()
	}
}
