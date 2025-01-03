package conf

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const (
	PROD = "production"
	DEV  = "development"
)

var (
	mu     sync.RWMutex
	config *Config
)

// Config 配置结构
type Config struct {
	RootURL     string         `mapstructure:"ROOT_URL"`
	Port        int            `mapstructure:"PORT"`
	Mode        string         `mapstructure:"MODE"`
	StaticPath  string         `mapstructure:"STATIC_PATH"`
	DataPath    string         `mapstructure:"DATA_PATH"`
	ApiVersion  string         `mapstructure:"API_VERSION"`
	MaxBodySize int64          `mapstructure:"MAX_BODY_SIZE"`
	Database    DatabaseConfig `mapstructure:"database"`
	JWT         JWTConfig      `mapstructure:"jwt"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `mapstructure:"HOST"`
	Port     string `mapstructure:"PORT"`
	Name     string `mapstructure:"NAME"`
	User     string `mapstructure:"USER"`
	Password string `mapstructure:"PASSWORD"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	AccessTokenMaxAge  time.Duration `mapstructure:"ACCESS_TOKEN_MAX_AGE"`
	RefreshTokenMaxAge time.Duration `mapstructure:"REFRESH_TOKEN_MAX_AGE"`
}

const (
	defaultSecretKey = "chang'emoonadsfwerf"
)

func init() {
	log.Println("初始化配置...")

	// 设置默认值
	viper.SetDefault("ROOT_URL", "http://localhost")
	viper.SetDefault("PORT", 8000)
	viper.SetDefault("MODE", "development")
	viper.SetDefault("STATIC_PATH", "/static")
	viper.SetDefault("DATA_PATH", filepath.Join("bin", "data"))
	viper.SetDefault("API_VERSION", "1.0")
	viper.SetDefault("MAX_BODY_SIZE", 10*1024*1024)

	viper.SetDefault("database.HOST", "localhost")
	viper.SetDefault("database.PORT", "5432")
	viper.SetDefault("database.NAME", "name")
	viper.SetDefault("database.USER", "postgres")

	viper.SetDefault("jwt.SECRET_KEY", defaultSecretKey)
	viper.SetDefault("jwt.ACCESS_TOKEN_MAX_AGE", 60*time.Minute)
	viper.SetDefault("jwt.REFRESH_TOKEN_MAX_AGE", 7*24*60*time.Minute)
	viper.SetDefault("jwt.PRIVATE_KEY", filepath.Join("bin", "config", "rsa_private_key.pem"))
	viper.SetDefault("jwt.PUBLIC_KEY", filepath.Join("bin", "config", "rsa_public_key.pem"))

	// 设置配置文件
	viper.SetConfigName("name")
	viper.SetConfigType("toml")
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

	// 初始化配置
	refreshConfig()

	// 监听配置文件变化
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed: %s", e.Name)
		refreshConfig()
	})
	viper.WatchConfig()
}

// refreshConfig 刷新配置
func refreshConfig() {
	mu.Lock()
	defer mu.Unlock()

	newConfig := &Config{}
	if err := viper.Unmarshal(newConfig); err != nil {
		log.Printf("Error unmarshaling config: %s", err)
		return
	}

	// 验证配置
	if err := validateConfig(newConfig); err != nil {
		log.Printf("Error validating config: %s", err)
		return
	}

	config = newConfig
	printConfig(config)
}

// validateConfig 验证配置的正确性
func validateConfig(cfg *Config) error {

	return nil
}

// GetConfig 返回配置实例
func GetConfig() *Config {
	mu.RLock()
	defer mu.RUnlock()
	return config
}

// Set 设置配置项
func Set(key string, value interface{}) {
	viper.Set(key, value)
	refreshConfig()
}

// Get 获取指定键的配置值
func Get(key string) interface{} {
	return viper.Get(key)
}

// GetString 获取字符串配置值
func GetString(key string) string {
	return viper.GetString(key)
}

// GetInt 获取整数配置值
func GetInt(key string) int {
	return viper.GetInt(key)
}

// GetBool 获取布尔配置值
func GetBool(key string) bool {
	return viper.GetBool(key)
}

// printConfig 打印配置信息
func printConfig(c *Config) {
	fmt.Println("===================配置文件================")
	fmt.Printf("ROOT_URL=%s\n", c.RootURL)
	fmt.Printf("PORT=%d\n", c.Port)
	fmt.Printf("MODE=%s\n", c.Mode)
	fmt.Printf("STATIC_PATH=%s\n", c.StaticPath)
	fmt.Printf("DATA_PATH=%s\n", c.DataPath)
	fmt.Printf("API_VERSION=%s\n", c.ApiVersion)

	fmt.Println("[DATABASE]")
	fmt.Printf("HOST=%s\n", c.Database.Host)
	fmt.Printf("NAME=%s\n", c.Database.Name)
	fmt.Printf("USER=%s\n", c.Database.User)

	fmt.Println("[JWT]")
	fmt.Printf("ACCESS_TOKEN_MAX_AGE=%d\n", c.JWT.AccessTokenMaxAge)
	fmt.Printf("REFRESH_TOKEN_MAX_AGE=%d\n", c.JWT.RefreshTokenMaxAge)

	fmt.Println("==========================================")
}
