package depts

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sshfs-admin/pkg/auth"
)

func SetupRouter(apiRouter *gin.RouterGroup) {
	route := apiRouter.Group("/depts")
	route.Use(auth.LoginRequiredMiddleware)
	route.GET("/", routeListDepts)
	route.GET("/count", routeCountDepts)
}

func routeCountDepts(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{})
}

func routeListDepts(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{})
}
