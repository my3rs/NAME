package conf

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

var configFile = "./name.conf"

type tomlConfig struct {
	Name       string `toml:"app_name"`
	Url        string `toml:ROOT_URL`
	Port       int
	Mode       Env
	StaticPath string       `toml:"STATIC_PATH"`
	DB         tomlDatabase `toml:"database"`
	JWT        tomlJWT      `toml:"jwt"`
	ApiVersion string       `toml:API_VERSION`
}

type tomlDatabase struct {
	Host     string
	Name     string
	User     string
	Password string
}

type tomlJWT struct {
	secretKey          string `toml:"SECRET_KEY"`
	AccessTokenMaxAge  int    `toml:"ACCESS_TOKEN_MAX_AGE"`
	RefreshTokenMaxAge int    `toml:"REFRESH_TOKEN_MAX_AGE"`
}

func (t *tomlJWT) SecretKey() string {
	if t.secretKey == "" {
		return defaultSecretKey
	}
	return t.secretKey
}

const (
	PROD             Env    = "production"
	DEV              Env    = "development"
	defaultSecretKey string = "chang'emoonadsfwerf"
)

type Env string

var (
	conf        *tomlConfig
	ProjectName = "NAME"
	Version     = "0.0.1"
)

func (c *tomlConfig) Print() {
	fmt.Println("===================配置文件================")
	fmt.Printf("NAME=%s\n", c.Name)
	fmt.Printf("PORT=%d\n", c.Port)
	fmt.Printf("MODE=%s\n", c.Mode)
	fmt.Printf("STATIC_PATH=%s\n", c.StaticPath)

	fmt.Println("[DATABASE]")
	fmt.Printf("HOST=%s\n", c.DB.Host)
	fmt.Printf("NAME=%s\n", c.DB.Name)
	fmt.Printf("USER=%s\n", c.DB.User)

	fmt.Println("[JWT]")
	fmt.Printf("SECRET_KEY=%s\n", c.JWT.SecretKey())
	fmt.Printf("ACCESS_TOKEN_MAX_AGE=%d\n", c.JWT.AccessTokenMaxAge)
	fmt.Printf("REFRESH_TOKEN_MAX_AGE=%d\n", c.JWT.RefreshTokenMaxAge)

	fmt.Println("==========================================")
}

// Config 单例模式
func Config() *tomlConfig {
	return GetConfig()
}

func initConfig() *tomlConfig {
	if _, err := toml.DecodeFile(configFile, &conf); err != nil {
		panic(err)
	}

	conf.Print()

	return conf
}

func GetConfig() *tomlConfig {
	if conf != nil {
		return conf
	} else {
		return initConfig()
	}
}

func Init(path string) *tomlConfig {
	configFile = path

	return GetConfig()
}
