package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"sshfs-admin/pkg/auth"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	store := memstore.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("auth", store))
	router.Static("/admin", "./admin-html")

	apiRouter := router.Group("/api")
	auth.SetupRouter(apiRouter)

	return router
}

func main() {
	engine := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	engine.Run(":8000")
}
