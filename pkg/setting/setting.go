package setting

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/go-ini/ini"
)

var (
	Cfg *ini.File

	RunMode string

	HttpPort       int
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int

	PageSize int

	Database struct {
		Type        string
		User        string
		Password    string
		Host        string
		Name        string
		TablePrefix string
	}

	JwtSecret      string
	JwtExpireHours int

	// 图片上传设置
	UploadPath   string
	ServerDomain string
)

func init() {
	// 尝试多个可能的路径
	configPaths := []string{
		"configs/app.ini",       // 相对于exe的同级configs目录
		"../../configs/app.ini", // 原来的路径
		"../configs/app.ini",    // 其他可能的路径
	}

	var loaded bool
	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			Cfg, err = ini.Load(path)
			if err == nil {
				log.Printf("成功加载配置文件: %s", path)
				loaded = true
				break
			}
		}
	}

	if !loaded {
		log.Fatalf("无法找到配置文件，尝试的路径: %v", configPaths)
	}

	LoadBase()
	LoadServer()
	LoadApp()
	LoadDatabase()
	LoadJWT()
	LoadUpload()
}

// LoadBase reads base configurations from the ini file
func LoadBase() {
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
}

func LoadServer() {
	sec, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatalf("Fail to get section 'server': %v", err)
	}

	HttpPort = sec.Key("HTTP_PORT").MustInt(8000)
	ReadTimeout = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60)) * time.Second
	WriteTimeout = time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second
	MaxHeaderBytes = sec.Key("Max_Header_Bytes").MustInt(1 << 20) // 1 MB
}

func LoadApp() {
	sec, err := Cfg.GetSection("app")
	if err != nil {
		log.Fatalf("Fail to get section 'app': %v", err)
	}

	PageSize = sec.Key("PAGE_SIZE").MustInt(10)
}

func LoadDatabase() {
	sec, err := Cfg.GetSection("database")
	if err != nil {
		log.Fatalf("Fail to get section 'database': %v", err)
	}

	Database.Type = sec.Key("TYPE").MustString("mysql")
	Database.User = sec.Key("USER").MustString("root")
	Database.Password = sec.Key("PASSWORD").MustString("1234")
	Database.Host = sec.Key("HOST").MustString("localhost:3306")
	Database.Name = sec.Key("NAME").MustString("botanical")
	Database.TablePrefix = sec.Key("TABLE_PREFIX").MustString("botanical_")
}

func LoadJWT() {
	sec, err := Cfg.GetSection("jwt")
	if err != nil {
		log.Fatalf("Fail to get section 'jwt': %v", err)
	}

	JwtSecret = sec.Key("JWT_SECRET").MustString("!@)*#)!@U#@*!@!)")
	JwtExpireHours = sec.Key("JWT_EXPIRE_HOURS").MustInt(24)
}

// LoadUpload 读取上传相关配置
func LoadUpload() {
	sec, err := Cfg.GetSection("upload")
	if err != nil {
		log.Fatalf("无法获取上传配置: %v", err)
	}

	UploadPath = sec.Key("UPLOAD_PATH").MustString("./uploads")
	ServerDomain = sec.Key("SERVER_DOMAIN").MustString("http://localhost:8000")

	// 确保上传目录存在
	if err := os.MkdirAll(filepath.Join(UploadPath, "avatars"), 0755); err != nil {
		log.Fatalf("创建头像上传目录失败: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(UploadPath, "pictures"), 0755); err != nil {
		log.Fatalf("创建图片上传目录失败: %v", err)
	}
}
