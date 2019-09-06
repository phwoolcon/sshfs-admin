package base

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/exec"
	"sshfs-admin/pkg/session"
	"strings"
)

const ConfigFile string = "/data/config.json"

var config Config
var Version string

type Config struct {
	loaded    bool
	HashSalt  string `json:"hash_salt"`
	SshfsHost string `json:"sshfs_host"`
	SshfsPort string `json:"sshfs_port"`
}

func GetConfig() Config {
	if !config.loaded {
		loadConfig()
	}
	return config
}

func IsFile(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func loadConfig() {
	configFile, err := os.Open(ConfigFile)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Loaded config " + ConfigFile)
	config.loaded = true
	fmt.Println(config)
}

func LocalExec(command string, arg ...string) (result []string) {
	fmt.Println("exec: ", command, arg)
	cmd := exec.Command(command, arg...)
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	output := strings.TrimSpace(string(cmdOutput.Bytes()))
	result = strings.Split(output, "\n")
	return result
}

func Route404(context *gin.Context) {
	if strings.HasPrefix(context.Request.RequestURI, "/api/") {
		context.JSON(http.StatusNotFound, gin.H{"error": "404 not found"})
		return
	}
	context.HTML(http.StatusNotFound, "404.html", nil)
}

func Session() gin.HandlerFunc {
	return sessions.Sessions("auth", session.NewFileStore("/data/session", []byte("secret")))
}
