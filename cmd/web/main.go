package main

import (
	"github.com/gin-gonic/gin"
	"path/filepath"
	"sshfs-admin/pkg/auth"
	"sshfs-admin/pkg/base"
	"sshfs-admin/pkg/depts"
	"sshfs-admin/pkg/system"
	"sshfs-admin/pkg/users"
	"strings"
)

var GinMode = ""
var Version = "dev"

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

func setupRouter(router *gin.Engine) {
	setupStaticRouter(router)
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

func setupStaticRouter(router *gin.Engine) {
	staticRouter := router.Group("/admin")
	staticRouter.Use(func(context *gin.Context) {
		filePath := context.Request.URL.Path
		if strings.HasSuffix(filePath, "/") || strings.ContainsRune(filepath.Base(filePath), '.') {
			context.Next()
			return
		}
		context.Params[0].Value += ".html"
		context.Request.URL.Path += ".html"
		context.Next()
	})
	staticRouter.Static("", "./web/admin")
}
