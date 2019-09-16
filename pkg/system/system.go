package system

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"sshfs-admin/pkg/auth"
	"sshfs-admin/pkg/base"
	"sshfs-admin/pkg/sshfs"
	"strconv"
	"strings"
)

func SetupRouter(apiRouter *gin.RouterGroup) {
	route := apiRouter.Group("/system")
	route.Use(auth.LoginRequiredMiddleware)
	route.GET("settings", routeGetConfig)
	route.GET("/status", routeStatus)
	route.POST("/settings", routeSaveSettings)
}

func routeGetConfig(context *gin.Context) {
	var config map[string]string
	configJson, _ := json.Marshal(base.GetConfig())
	_ = json.Unmarshal(configJson, &config)
	delete(config, "hash_salt")
	response := gin.H{"config": config}
	context.JSON(http.StatusOK, response)
}

func routeSaveSettings(context *gin.Context) {
	config := base.GetConfig()
	sshfsHost := context.PostForm("sshfs_host")
	sshfsPort := context.PostForm("sshfs_port")
	httpsHost := context.PostForm("https_host")
	httpsPort := context.PostForm("https_port")
	if !verifyHost(sshfsHost) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sshfs host"})
		return
	}
	sshfsPortNum, err := strconv.Atoi(sshfsPort)
	if err != nil || sshfsPortNum < 1 || sshfsPortNum > 65535 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Sshfs port must be an integer between 1 and 65535"})
		return
	}
	if len(httpsHost) > 0 && !verifyHost(httpsHost) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid https host"})
		return
	}
	httpsPortNum, err := strconv.Atoi(httpsPort)
	if len(httpsPort) > 0 && (err != nil || httpsPortNum < 1 || httpsPortNum > 65535) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Https port must be an integer between 1 and 65535"})
		return
	}
	config.SshfsHost = sshfsHost
	config.SshfsPort = sshfsPort
	config.HttpsHost = httpsHost
	config.HttpsPort = httpsPort
	base.SaveConfig(config)
	context.JSON(http.StatusOK, gin.H{"message": "Changes saved"})
}

func routeStatus(context *gin.Context) {
	usages := strings.Fields(sshfs.GetDiskUsage()[0])
	status := make(map[string]string)
	if len(usages) >= 3 {
		status["used"] = usages[0]
		status["free"] = usages[1]
		status["total"] = usages[2]
		status["free_percent"] = usages[3]
	}
	context.JSON(http.StatusOK, gin.H{"status": status})
}

func verifyHost(name string) bool {
	validName := regexp.MustCompile(`^[\w][\w.\-]+$`)
	return validName.MatchString(name)
}
