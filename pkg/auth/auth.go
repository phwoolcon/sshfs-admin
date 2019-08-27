package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strconv"
)

const usersFile string = "/data/users.json"

type User struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	SessionTtl int    `json:"session_ttl"`
}

type Users map[string]User

func LoginRequiredMiddleware(context *gin.Context) {
	user, err := getSessionUser(context);
	if err != nil {
		context.Abort()
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Login required"})
		return
	}
	context.Set("user", user)
	context.Next()
}

func SetupRouter(apiRouter *gin.RouterGroup) {
	public := apiRouter.Group("/auth")
	public.GET("/status", routeStatus)
	public.GET("/logout", routeLogout)
	public.POST("/login", routeLogin)

	private := apiRouter.Group("/auth")
	private.Use(LoginRequiredMiddleware)
	private.GET("/settings", routeGetSettings)
	private.POST("/settings", routeSaveSettings)
	private.POST("/change-pass", routeChangePassword)
}

func getSessionUser(context *gin.Context) (user User, err error) {
	session := sessions.Default(context)
	username := session.Get("username")
	return findUser(fmt.Sprintf("%v", username))
}

func findUser(username string) (user User, err error) {
	users := loadUsers()
	user, ok := users[username]
	if ok {
		return user, nil
	}
	return User{}, errors.New("Attempt to login as non-existing user: " + username)
}

func loadUsers() (users Users) {
	usersFile, err := os.Open(usersFile)
	defer usersFile.Close()
	if err != nil {
		fmt.Println(err)
		return Users{}
	}
	err = json.NewDecoder(usersFile).Decode(&users)
	if err != nil {
		fmt.Println(err)
		return Users{}
	}
	return users
}

func routeChangePassword(context *gin.Context) {
	oldPassword := context.PostForm("old_password")
	newPassword := context.PostForm("new_password")
	user := context.MustGet("user").(User)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		fmt.Println("Attempt to change password failed for user: " + user.Username)
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credential"})
		return
	}
	fmt.Println("Password changed successfully for user: " + user.Username)
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 5)
	user.Password = string(newHashedPassword)
	saveUser(user)
	context.JSON(http.StatusOK, gin.H{"message": "Password changed"})
}

func routeGetSettings(context *gin.Context) {
	user := context.MustGet("user").(User)
	context.JSON(http.StatusOK, gin.H{"session_ttl": user.SessionTtl})
}

func routeLogin(context *gin.Context) {
	username := context.PostForm("username")
	password := context.PostForm("password")

	user, err := findUser(username)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println("User " + username + " login failed")
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credential"})
		return
	}
	session := sessions.Default(context)
	session.Options(sessions.Options{
		Path:     "/",
		MaxAge:   user.SessionTtl,
		HttpOnly: true,
	})
	session.Set("username", user.Username)
	session.Save()
	fmt.Println("User " + username + " logged in successfully with ttl " + strconv.Itoa(user.SessionTtl))
	context.JSON(http.StatusOK, gin.H{"username": username})
}

func routeLogout(context *gin.Context) {
	session := sessions.Default(context)
	session.Options(sessions.Options{
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	session.Set("username", nil)
	session.Save()
	context.JSON(http.StatusOK, gin.H{"username": nil})
}

func routeSaveSettings(context *gin.Context) {
	sessionTtl, err := strconv.Atoi(context.PostForm("session_ttl"))
	if err != nil || sessionTtl < 60 || sessionTtl > 86400*7 {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Session TTL should be between 60 and 604800"})
		return
	}
	user := context.MustGet("user").(User)
	fmt.Println("Settings changed successfully for user: " + user.Username)
	user.SessionTtl = sessionTtl
	saveUser(user)
	context.JSON(http.StatusOK, gin.H{"message": "Settings changed"})
}

func routeStatus(context *gin.Context) {
	user, _ := getSessionUser(context)
	context.JSON(http.StatusOK, gin.H{"username": user.Username})
}

func saveUser(user User) {
	users := loadUsers()
	usersFile, err := os.OpenFile(usersFile, os.O_WRONLY|os.O_CREATE, 0644)
	defer usersFile.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	jsonEncoder := json.NewEncoder(usersFile)
	jsonEncoder.SetIndent("", "    ")
	users[user.Username] = user
	err = jsonEncoder.Encode(users);
	if err != nil {
		fmt.Println(err)
		return
	}
}
