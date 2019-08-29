package depts

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sshfs-admin/pkg/auth"
	"sshfs-admin/pkg/sshfs"
)

func SetupRouter(apiRouter *gin.RouterGroup) {
	route := apiRouter.Group("/depts")
	route.Use(auth.LoginRequiredMiddleware)
	route.GET("", routeListDepts)
	route.GET("/count", routeCountDepts)
}

func routeCountDepts(context *gin.Context) {
	departments := sshfs.GetDepartments()
	context.JSON(http.StatusOK, gin.H{"count": len(departments)})
}

func routeListDepts(context *gin.Context) {
	departments := sshfs.GetDepartments()
	context.JSON(http.StatusOK, gin.H{"depts": departments})
}
