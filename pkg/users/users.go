package users

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"sshfs-admin/pkg/auth"
	"sshfs-admin/pkg/sshfs"
)

func SetupRouter(apiRouter *gin.RouterGroup) {
	route := apiRouter.Group("/users")
	route.Use(auth.LoginRequiredMiddleware)
	route.GET("", routeList)
	route.GET("/count", routeCount)
	route.POST("/create", routeCreate)
}

func routeCreate(context *gin.Context) {
	name := context.PostForm("name")
	department := context.PostForm("dept")
	validName := regexp.MustCompile(`^[A-Za-z][\w.\-]{2,14}$`)
	if !validName.MatchString(name) {
		fmt.Println("Invalid user name: " + name)
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "User name must begin with a letter, and be between 3 and 15 characters of \"A-Za-z0-9.-_\"",
		})
		return
	}
	result := sshfs.CreateUser(name, department)
	if result[0] != "ok" {
		context.JSON(http.StatusBadRequest, gin.H{"error": result[0]})
		return
	}
	context.JSON(http.StatusOK, gin.H{})
}

func routeCount(context *gin.Context) {
	users := sshfs.GetUsers()
	context.JSON(http.StatusOK, gin.H{"count": len(users)})
}

func routeList(context *gin.Context) {
	users := sshfs.GetUsers()
	context.JSON(http.StatusOK, gin.H{"users": users})
}
