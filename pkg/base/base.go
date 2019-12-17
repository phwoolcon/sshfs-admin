package base

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/phwoolcon/gin-utils/session"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const ConfigFile string = "/data/config.json"
const SshfsHostFile string = "/data/sshfs_host"

var config Config
var Version string

type Config struct {
	loaded      bool
	HashSalt    string `json:"hash_salt"`
	SshfsHost   string `json:"sshfs_host"`
	SshfsPort   string `json:"sshfs_port"`
	RawHttpPort string `json:"raw_http_port"`
	HttpsHost   string `json:"https_host"`
	HttpsPort   string `json:"https_port"`
}

func GetConfig() Config {
	if !config.loaded {
		loadConfig()
		initHashSalt()
	}
	return config
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// Get preferred outbound ip of this machine https://stackoverflow.com/a/37382208/802646
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func HasTlsCert(certFile, keyFile string) bool {
	if !IsFile(certFile) || !IsFile(keyFile) {
		return false
	}
	_, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return false
	}
	return true
}

func initHashSalt() {
	if len(strings.TrimSpace(config.HashSalt)) >= 16 {
		return
	}
	salt := make([]byte, 16)
	saltBase64 := make([]byte, base64.RawURLEncoding.EncodedLen(len(salt)))
	_, _ = rand.Read(salt)
	base64.RawURLEncoding.Encode(saltBase64, salt)
	config.HashSalt = string(saltBase64)[0:16]
	fmt.Println("Generated hash salt: " + config.HashSalt)
	SaveConfig(config)
}

func IsFile(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func IsTls(context *gin.Context) bool {
	request := context.Request
	return request.URL.Scheme == "https" || request.TLS != nil || context.GetHeader("x-forwarded-proto") == "https"
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
	hostConfig := LocalExec("cat", SshfsHostFile)
	config.SshfsHost = hostConfig[0]
	config.SshfsPort = hostConfig[1]
	config.RawHttpPort = hostConfig[2]
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

func RedirectToHttpsMiddleware(context *gin.Context) {
	if IsTls(context) {
		context.Next()
		return
	}
	GetConfig()
	url := *context.Request.URL
	url.Scheme = "https"
	url.Host = config.HttpsHost
	if config.HttpsPort != "443" {
		url.Host += ":" + config.HttpsPort
	}
	context.Redirect(http.StatusFound, url.String())
	context.Abort()
}

func Route404(context *gin.Context) {
	if strings.HasPrefix(context.Request.RequestURI, "/api/") {
		context.JSON(http.StatusNotFound, gin.H{"error": "404 not found"})
		return
	}
	context.HTML(http.StatusNotFound, "404.html", nil)
}

func SaveConfig(newConfig Config) {
	configFile, err := os.OpenFile(ConfigFile, os.O_WRONLY|os.O_CREATE, 0644)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	jsonEncoder := json.NewEncoder(configFile)
	jsonEncoder.SetIndent("", "    ")
	err = jsonEncoder.Encode(newConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	loadConfig()
}

func SessionMiddleware() gin.HandlerFunc {
	return sessions.Sessions("auth", session.NewFileStore("/data/session", []byte("secret")))
}

func SshfsRootAccess() []string {
	return LocalExec("./scripts/sshfs_root_access")
}
