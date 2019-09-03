package users

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/gin-gonic/gin"
	"hash/crc32"
	"net/http"
	"regexp"
	"sshfs-admin/pkg/auth"
	"sshfs-admin/pkg/base"
	"sshfs-admin/pkg/sshfs"
	"strings"
)

func SetupApiRouter(apiRouter *gin.RouterGroup) {
	route := apiRouter.Group("/users")
	route.Use(auth.LoginRequiredMiddleware)
	route.GET("", routeList)
	route.GET("/count", routeCount)
	route.GET("/details", routeDetails)
	route.POST("/create", routeCreate)
}

func SetupFrontRouter(router *gin.Engine) {
	router.GET("/download/:token", frontRouteDownloadPage)
}

func frontRouteDownloadPage(context *gin.Context) {
	token := strings.Split(context.Param("token"), "~")
	if len(token) != 2 {
		base.Route404(context)
		return
	}
	name := token[0]
	hash := token[1]
	expectedHash := getUserHash(name)
	if !verifyName(name) || hash != expectedHash {
		base.Route404(context)
		return
	}
	context.HTML(http.StatusOK, "download.html", nil)
}

func getUserHash(name string) (hash string) {
	sshKeyHash := base.LocalExec("./scripts/sshkey_md5", name)[0]
	crc32Checksum := crc32.ChecksumIEEE([]byte(name + "|" + base.GetConfig().SecretKey + sshKeyHash))
	crc32Byte := make([]byte, 4)
	crc32Base64 := make([]byte, base64.RawURLEncoding.EncodedLen(len(crc32Byte)))
	binary.LittleEndian.PutUint32(crc32Byte, crc32Checksum)
	base64.RawURLEncoding.Encode(crc32Base64, []byte(crc32Byte))
	return string(crc32Base64)
}

func routeCreate(context *gin.Context) {
	name := context.PostForm("name")
	department := context.PostForm("dept")
	if !verifyName(name) {
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

func routeDetails(context *gin.Context) {
	user := context.Query("user")
	if !verifyName(user) {
		base.Route404(context)
		return
	}
	department := sshfs.GetUserDepartments(user)[0]
	context.JSON(http.StatusOK, gin.H{"dept": department, "token": user + "~" + getUserHash(user)})
}

func routeList(context *gin.Context) {
	users := sshfs.GetUsers()
	context.JSON(http.StatusOK, gin.H{"users": users})
}

func verifyName(name string) bool {
	validName := regexp.MustCompile(`^[A-Za-z][\w.\-]{2,14}$`)
	return validName.MatchString(name)
}
