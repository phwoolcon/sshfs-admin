package main

import (
	"github.com/gin-gonic/gin"
	"sshfs-admin/pkg/auth"
	"sshfs-admin/pkg/base"
	"sshfs-admin/pkg/depts"
	"sshfs-admin/pkg/system"
	"sshfs-admin/pkg/users"
)

var GinMode string
var Version string = "dev"

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Static("/admin", "./web/admin")
	router.LoadHTMLFiles("./web/admin/404.html", "./web/admin/download.html")
	router.NoRoute(base.Route404)

	apiRouter := router.Group("/api")
	apiRouter.Use(base.Session())

	auth.SetupRouter(apiRouter)
	depts.SetupRouter(apiRouter)
	users.SetupApiRouter(apiRouter)
	users.SetupFrontRouter(router)
	system.SetupRouter(apiRouter)

	return router
}

func main() {
	if GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	accessStatus := base.SshfsRootAccess()
	if accessStatus[0] != "ok" {
		panic(accessStatus)
	}
	base.Version = Version
	engine := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	engine.Run(":8000")
}
