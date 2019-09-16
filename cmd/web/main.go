package main

import (
	"github.com/gin-gonic/gin"
	"sshfs-admin/pkg/auth"
	"sshfs-admin/pkg/base"
	"sshfs-admin/pkg/depts"
	"sshfs-admin/pkg/system"
	"sshfs-admin/pkg/users"
)

var GinMode = ""
var Version = "dev"

func setupRouter(router *gin.Engine) {
	router.Static("/admin", "./web/admin")
	router.LoadHTMLFiles("./web/admin/404.html", "./web/admin/download.html")
	router.NoRoute(base.Route404)

	apiRouter := router.Group("/api")
	apiRouter.Use(base.SessionMiddleware())

	auth.SetupRouter(apiRouter)
	depts.SetupRouter(apiRouter)
	users.SetupApiRouter(apiRouter)
	users.SetupFrontRouter(router)
	system.SetupRouter(apiRouter)
}

func main() {
	gin.SetMode(GinMode)
	accessStatus := base.SshfsRootAccess()
	if accessStatus[0] != "ok" {
		panic(accessStatus)
	}
	base.Version = Version

	engine := gin.Default()

	if len(base.GetConfig().HttpsHost) > 0 {
		engine.Use(base.RedirectToHttpsMiddleware)
	}
	setupRouter(engine)
	tlsCertFile := "/data/tls/cert"
	tlsKeyFile := "/data/tls/key"
	hasTlsCert := base.HasTlsCert(tlsCertFile, tlsKeyFile)
	if hasTlsCert {
		go engine.RunTLS(":8443", tlsCertFile, tlsKeyFile)
	}
	engine.Run(":8000")
}
