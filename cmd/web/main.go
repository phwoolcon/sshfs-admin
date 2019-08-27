package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"net/http"
	"sshfs-admin/pkg/auth"
	"sshfs-admin/pkg/depts"
	"strings"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Static("/admin", "./admin-html")
	router.LoadHTMLFiles("admin-html/404.html")
	router.NoRoute(func(context *gin.Context) {
		if strings.HasPrefix(context.Request.RequestURI, "/admin/") {
			context.HTML(http.StatusNotFound, "404.html", nil)
			return
		}
		context.JSON(http.StatusNotFound, gin.H{"error": "404 not found"})
	})

	apiRouter := router.Group("/api")
	apiRouter.Use(sessions.Sessions("auth", memstore.NewStore([]byte("secret"))))

	auth.SetupRouter(apiRouter)
	depts.SetupRouter(apiRouter)

	return router
}

func main() {
	engine := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	engine.Run(":8000")
}
