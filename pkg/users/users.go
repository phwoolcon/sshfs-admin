package users

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sshfs-admin/pkg/auth"
	"sshfs-admin/pkg/sshfs"
)

func SetupRouter(apiRouter *gin.RouterGroup) {
	route := apiRouter.Group("/users")
	route.Use(auth.LoginRequiredMiddleware)
	route.GET("", routeListUsers)
	route.GET("/count", routeCountUsers)
}

func routeCountUsers(context *gin.Context) {
	users := sshfs.GetUsers()
	context.JSON(http.StatusOK, gin.H{"count": len(users)})
}

func routeListUsers(context *gin.Context) {
	users := sshfs.GetUsers()
	context.JSON(http.StatusOK, gin.H{"users": users})
}
