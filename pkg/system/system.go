package system

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"sshfs-admin/pkg/auth"
	"sshfs-admin/pkg/base"
	"strconv"
)

func SetupRouter(apiRouter *gin.RouterGroup) {
	route := apiRouter.Group("/system")
	route.Use(auth.LoginRequiredMiddleware)
	route.GET("", routeGetConfig)
	route.POST("/sshfs", routeSaveSshfsConfig)
}

func routeGetConfig(context *gin.Context) {
	var config map[string]string
	configJson, _ := json.Marshal(base.GetConfig())
	_ = json.Unmarshal(configJson, &config)
	delete(config, "hash_salt")
	response := gin.H{"config": config}
	context.JSON(http.StatusOK, response)
}

func routeSaveSshfsConfig(context *gin.Context) {
	config := base.GetConfig()
	host := context.PostForm("sshfs_host")
	port := context.PostForm("sshfs_port")
	if !verifyHost(host) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host"})
		return
	}
	portNum, err := strconv.Atoi(port)
	if err != nil || portNum < 1 || portNum > 65535 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Port must be a integer between 1 and 65535"})
		return
	}
	config.SshfsHost = host
	config.SshfsPort = port
	base.SaveConfig(config)
	context.JSON(http.StatusOK, gin.H{"message": "Changes saved"})
}

func verifyHost(name string) bool {
	validName := regexp.MustCompile(`^[\w][\w.\-]+$`)
	return validName.MatchString(name)
}