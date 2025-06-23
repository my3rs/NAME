package conf

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	PROD = "production"
	DEV  = "development"
)

const (
	DATABASE_DRIVER_SQLITE   = "sqlite"
	DATABASE_DRIVER_POSTGRES = "postgres"
)

var (
	config *Config
)

// Config 配置结构
type Config struct {
	Host        string         `mapstructure:"HOST"`
	Port        int            `mapstructure:"PORT"`
	Mode        string         `mapstructure:"MODE"`
	AssetsPath  string         `mapstructure:"ASSETS_PATH"`
	DataPath    string         `mapstructure:"DATA_PATH"`
	UploadsPath string         `mapstructure:"UPLOADS_PATH"`
	MaxBodySize int64          `mapstructure:"MAX_BODY_SIZE"`
	Database    DatabaseConfig `mapstructure:"database"`
	JWT         JWTConfig      `mapstructure:"jwt"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver   string `mapstructure:"DRIVER"`
	DataPath string `mapstructure:"DATA_PATH"`
	FileName string `mapstructure:"FILE_NAME"`
	Host     string `mapstructure:"HOST"`
	Port     string `mapstructure:"PORT"`
	Name     string `mapstructure:"NAME"`
	User     string `mapstructure:"USER"`
	Password string `mapstructure:"PASSWORD"`
	SSLMode  string `mapstructure:"SSL_MODE"`
	TimeZone string `mapstructure:"TIME_ZONE"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret             string        `mapstructure:"SECRET_KEY"`
	Issuer             string        `mapstructure:"ISSUER"`
	AccessTokenMaxAge  time.Duration `mapstructure:"ACCESS_TOKEN_MAX_AGE"`
	RefreshTokenMaxAge time.Duration `mapstructure:"REFRESH_TOKEN_MAX_AGE"`
}

func init() {
	log.Println("初始化配置...")

	// 设置默认值
	viper.SetDefault("HOST", "127.0.0.1")
	viper.SetDefault("PORT", 8000)
	viper.SetDefault("MODE", DEV)
	viper.SetDefault("ASSETS_PATH", "./web/assets")
	viper.SetDefault("DATA_PATH", filepath.Join("bin", "data"))
	viper.SetDefault("UPLOADS_PATH", "/uploads")
	viper.SetDefault("MAX_BODY_SIZE", 10*1024*1024)

	viper.SetDefault("database.DRIVER", DATABASE_DRIVER_SQLITE)
	viper.SetDefault("database.DATA_PATH", ".")
	viper.SetDefault("database.FILE_NAME", "name.db")
	viper.SetDefault("database.HOST", "localhost")
	viper.SetDefault("database.PORT", "5432")
	viper.SetDefault("database.NAME", "name")
	viper.SetDefault("database.USER", "postgres")
	viper.SetDefault("database.SSL_MODE", "disable")
	viper.SetDefault("database.TIME_ZONE", "Asia/Shanghai")

	viper.SetDefault("jwt.ACCESS_TOKEN_MAX_AGE", 60*time.Minute)
	viper.SetDefault("jwt.REFRESH_TOKEN_MAX_AGE", 7*24*60*time.Minute)
	viper.SetDefault("jwt.ISSUER", "NAME")

	// 设置配置文件
	viper.SetConfigName("name")
	viper.SetConfigType("toml")
	viper.AddConfigPath("$HOME/.name")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("./bin/config")
	viper.AddConfigPath(".")

	// 读取环境变量
	viper.AutomaticEnv()
	viper.SetEnvPrefix("NAME")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 尝试读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// 配置文件存在但读取失败
			panic(fmt.Errorf("fatal error reading config file: %s", err))
		}
		// 配置文件不存在，使用默认配置
		log.Println("No config file found, using defaults")
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Panic("unmarshal config err:", err)
		return
	}

	printConfig(config)
}

// GetConfig 返回配置实例
func GetConfig() *Config {
	return config
}

// printConfig 打印配置信息
func printConfig(c *Config) {
	fmt.Println("===================当前配置================")
	fmt.Printf("HOST=%s\n", c.Host)
	fmt.Printf("PORT=%d\n", c.Port)
	fmt.Printf("MODE=%s\n", c.Mode)
	fmt.Printf("ASSETS_PATH=%s\n", c.AssetsPath)
	fmt.Printf("DATA_PATH=%s\n", c.DataPath)
	fmt.Printf("UPLOADS_PATH=%s\n", c.UploadsPath)

	fmt.Println("[DATABASE]")
	fmt.Printf("DRIVER=%s\n", c.Database.Driver)
	fmt.Printf("DATA_PATH=%s\n", c.Database.DataPath)
	fmt.Printf("FILE_NAME=%s\n", c.Database.FileName)
	fmt.Printf("HOST=%s\n", c.Database.Host)
	fmt.Printf("NAME=%s\n", c.Database.Name)
	fmt.Printf("USER=%s\n", c.Database.User)
	fmt.Printf("SSL_MODE=%s\n", c.Database.SSLMode)
	fmt.Printf("TIME_ZONE=%s\n", c.Database.TimeZone)

	fmt.Println("[JWT]")
	fmt.Printf("SECRET=%s\n", c.JWT.Secret)
	fmt.Printf("ISSUER=%s\n", c.JWT.Issuer)
	fmt.Printf("ACCESS_TOKEN_MAX_AGE=%d\n", c.JWT.AccessTokenMaxAge)
	fmt.Printf("REFRESH_TOKEN_MAX_AGE=%d\n", c.JWT.RefreshTokenMaxAge)

	fmt.Println("==========================================")
}
