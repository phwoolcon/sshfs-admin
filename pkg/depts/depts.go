package depts

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"sshfs-admin/pkg/auth"
	"sshfs-admin/pkg/sshfs"
	"strings"
)

type DepartmentUsage struct {
	Name  string `json:"name"`
	Usage string `json:"usage"`
}

func SetupRouter(apiRouter *gin.RouterGroup) {
	route := apiRouter.Group("/depts")
	route.Use(auth.LoginRequiredMiddleware)
	route.GET("", routeList)
	route.GET("/options", routeOptions)
	route.GET("/count", routeCount)
	route.POST("/create", routeCreate)
}

func routeCreate(context *gin.Context) {
	validName := regexp.MustCompile(`^[A-Za-z][\w.\-]{2,14}$`)
	name := context.PostForm("name")
	if !validName.MatchString(name) {
		fmt.Println("Invalid department name: " + name)
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Department name must begin with a letter, and be between 3 and 15 characters of \"A-Za-z0-9.-_\"",
		})
		return
	}
	result := sshfs.CreateDepartment(name)
	if result[0] != "ok" {
		context.JSON(http.StatusBadRequest, gin.H{"error": result[0]})
		return
	}
	context.JSON(http.StatusOK, gin.H{})
}

func routeCount(context *gin.Context) {
	count := sshfs.GetDepartmentCount()[0]
	context.JSON(http.StatusOK, gin.H{"count": count})
}

func routeList(context *gin.Context) {
	departmentUsages := sshfs.GetDepartmentsWithUsages()
	departments := make([]DepartmentUsage, 0)
	for _, usage := range departmentUsages {
		usageInfo := strings.Fields(usage)
		if len(usageInfo) != 2 {
			continue
		}
		departments = append(departments, DepartmentUsage{Name: usageInfo[1], Usage: usageInfo[0]})
	}
	context.JSON(http.StatusOK, gin.H{"depts": departments})
}

func routeOptions(context *gin.Context) {
	options := make(map[string]string)
	for _, department := range sshfs.GetDepartments() {
		options[department] = department
	}
	context.JSON(http.StatusOK, gin.H{"options": options})
}
