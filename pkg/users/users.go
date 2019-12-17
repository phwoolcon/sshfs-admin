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
	"strconv"
	"strings"
	"time"
)

type UserUsage struct {
	Name  string `json:"name"`
	Usage string `json:"usage"`
}

func SetupApiRouter(apiRouter *gin.RouterGroup) {
	route := apiRouter.Group("/users")
	route.Use(auth.LoginRequiredMiddleware)
	route.GET("", routeList)
	route.GET("/count", routeCount)
	route.GET("/details", routeDetails)
	route.POST("/create", routeCreate)
	route.POST("/edit", routeEdit)
}

func SetupFrontRouter(router *gin.Engine) {
	router.GET("/download/:token", frontRouteDownloadPage)
	router.GET("/download/:token/config", frontRouteDownloadConfig)
	router.GET("/download/:token/key", frontRouteDownloadKey)
	router.GET("/api/users/has-key/:token", frontRouteHasKey)
	router.POST("/api/users/regenerate-key/:token", frontRouteRegenerateKey)
	router.POST("/api/users/front-create", frontRouteCreate)
}

func convertTokenToUserName(tokenString string) string {
	token := strings.Split(tokenString, "~")
	if len(token) != 2 {
		return ""
	}
	name := token[0]
	hash := token[1]
	expectedHash := getUserHash(name)
	if !verifyName(name) || hash != expectedHash {
		return ""
	}
	if !sshfs.UserExists(name) {
		return ""
	}
	return name
}

func frontRouteCreate(context *gin.Context) {
	context.Set("front-create", true)
	routeCreate(context)
}

func frontRouteDownloadConfig(context *gin.Context) {
	name := convertTokenToUserName(context.Param("token"))
	if name == "" {
		base.Route404(context)
		return
	}
	configTemplate := `HOST=%s
PORT=%s
USER=%s
DRIVE=Z

`
	app := base.GetConfig()
	context.Header("content-disposition", "attachment; filename=config.ini")
	context.String(http.StatusOK, configTemplate, app.SshfsHost, app.SshfsPort, name)
}

func frontRouteDownloadKey(context *gin.Context) {
	name := convertTokenToUserName(context.Param("token"))
	if name == "" {
		base.Route404(context)
		return
	}
	privateKey := base.LocalExec("./scripts/sshkey_download", name)
	if privateKey[0] == "" {
		base.Route404(context)
		return
	}
	context.Header("content-disposition", "attachment; filename=ssh.key")
	context.String(http.StatusOK, strings.Join(privateKey, "\n"))
}

func frontRouteDownloadPage(context *gin.Context) {
	name := convertTokenToUserName(context.Param("token"))
	if name == "" {
		base.Route404(context)
		return
	}
	context.HTML(http.StatusOK, "download.html", nil)
}

func frontRouteHasKey(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"result": hasPrivateKey(context.Param("token"))})
}

func frontRouteRegenerateKey(context *gin.Context) {
	name := convertTokenToUserName(context.Param("token"))
	if name == "" {
		base.Route404(context)
		return
	}
	result := sshfs.RegenerateKey(name)
	if result[0] != "ok" {
		context.JSON(http.StatusBadRequest, gin.H{"error": result[0]})
		return
	}
	context.JSON(http.StatusOK, gin.H{"new_token": getUserToken(name)})
}

func getUserHash(name string) (hash string) {
	sshKeySum := base.LocalExec("./scripts/sshkey_sum", name)[0]
	crc32Checksum := crc32.ChecksumIEEE([]byte(name + "|" + base.GetConfig().HashSalt + sshKeySum))
	crc32Byte := make([]byte, 4)
	crc32Base64 := make([]byte, base64.RawURLEncoding.EncodedLen(len(crc32Byte)))
	binary.LittleEndian.PutUint32(crc32Byte, crc32Checksum)
	base64.RawURLEncoding.Encode(crc32Base64, crc32Byte)
	return string(crc32Base64)
}

func getUserToken(name string) string {
	return name + "~" + getUserHash(name)
}

func hasPrivateKey(token string) bool {
	name := convertTokenToUserName(token)
	if name == "" {
		return false
	}
	return base.IsFile(fmt.Sprintf(`/data/tmp/%s.key`, name))
}

func routeCreate(context *gin.Context) {
	name := context.PostForm("name")
	department := context.PostForm("dept")
	if !verifyName(name) {
		fmt.Println("Invalid user name: " + name)
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "User name must begin with a letter, and be between 3 and 32 characters of \"A-Za-z0-9.-_\"",
		})
		return
	}
	if !verifyName(department) {
		fmt.Println("Invalid department name: " + department)
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid department",
		})
		return
	}
	if front, hasValue := context.Get("front-create"); hasValue && front.(bool) {
		t := context.PostForm("t")
		timestamp, _ := strconv.Atoi(t)
		if int64(timestamp) < time.Now().Unix()-600 {
			fmt.Println("Invalid t:", t)
			context.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid t",
			})
			return
		}
		hash := base.GetMD5Hash(name + department + string(t) + base.GetConfig().HashSalt)
		sign := context.PostForm("sign")
		if sign != hash {
			fmt.Println("Invalid sign:", sign, "expected:", hash)
			context.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid sign",
			})
			return
		}
		sshfs.CreateDepartment(department)
	}
	result := sshfs.CreateUser(name, department)
	if result[0] != "ok" {
		context.JSON(http.StatusBadRequest, gin.H{"error": result[0]})
		return
	}
	// TODO Generate https link
	linkScheme := "http"
	linkHost := base.GetConfig().SshfsHost
	linkPort := base.GetConfig().RawHttpPort
	if linkPort != "80" {
		linkHost += ":" + linkPort
	}
	context.JSON(http.StatusOK, gin.H{"usage": fmt.Sprintf(`Open this link to download client and config:
%s://%s/download/%s`, linkScheme, linkHost, getUserToken(name))})
}

func routeCount(context *gin.Context) {
	count := sshfs.GetUserCount()[0]
	context.JSON(http.StatusOK, gin.H{"count": count})
}

func routeDetails(context *gin.Context) {
	user := context.Query("user")
	if !verifyName(user) {
		base.Route404(context)
		return
	}
	if !sshfs.UserExists(user) {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "The user you are editing does not exist",
		})
		return
	}
	department := sshfs.GetUserDepartments(user)[0]
	context.JSON(http.StatusOK, gin.H{"dept": department, "token": getUserToken(user)})
}

func routeEdit(context *gin.Context) {
	name := context.PostForm("orig_name")
	newName := context.PostForm("name")
	newDepartment := context.PostForm("dept")
	if !verifyName(newName) || !verifyName(name) {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "User name must begin with a letter, and be between 3 and 32 characters of \"A-Za-z0-9.-_\"",
		})
		return
	}
	if !sshfs.UserExists(name) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "The user you are editing does not exist"})
		return
	}
	department := sshfs.GetUserDepartments(name)[0]
	if name != newName {
		renameResult := sshfs.RenameUser(name, newName)
		if renameResult[0] != "ok" {
			context.JSON(http.StatusBadRequest, gin.H{
				"error": renameResult[0],
			})
			return
		}
	}
	if department != newDepartment {
		renameResult := sshfs.UpdateUserDepartment(newName, newDepartment)
		if renameResult[0] != "ok" {
			context.JSON(http.StatusBadRequest, gin.H{
				"error": renameResult[0],
			})
			return
		}
	}
	context.JSON(http.StatusOK, gin.H{"user": newName, "message": "Changes saved"})
}

func routeList(context *gin.Context) {
	userUsages := sshfs.GetUsersWithUsages()
	users := make([]UserUsage, 0)
	for _, usage := range userUsages {
		usageInfo := strings.Fields(usage)
		if len(usageInfo) != 2 {
			continue
		}
		users = append(users, UserUsage{Name: usageInfo[1], Usage: usageInfo[0]})
	}
	context.JSON(http.StatusOK, gin.H{"users": users})
}

func verifyName(name string) bool {
	validName := regexp.MustCompile(`^[\w.\-]{3,32}$`)
	return validName.MatchString(name)
}
